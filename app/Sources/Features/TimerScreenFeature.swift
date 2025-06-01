//
//  TimerScreenFeature.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/04/19.
//
import ComposableArchitecture
import SwiftUI

@Reducer
struct TimerScreenFeature {
    struct State: Equatable {
        var timer: TimerFeature.State
        var evalModal: EvalModalFeature.State?
        var sessionId: UUID?

    }

    enum Action {
        case timer(TimerFeature.Action)
        case evalModal(EvalModalFeature.Action)

        case fetchCycleTapped

        case startNextRound
        
        case sessionStartResponse(Result<SessionResult, Error>)
        case roundStartResponse(Result<RoundResult, Error>)
        
        case showEvalModal
        case dismissEvalModal
    }
    
    @Dependency(\.sessionAPIClient) var sessionAPIClient
    
    var body: some ReducerOf<Self> {
        Scope(state: \.timer, action: \.timer){TimerFeature()}
        .ifLet(\.evalModal, action: \.evalModal) {EvalModalFeature()}
        
        Reduce {state, action in
            switch action {
                
            case .fetchCycleTapped:
                return .run { send in
                    let session = try await sessionAPIClient.startSession()
                    print("Start session success → \(session.id)")
                    await send(.sessionStartResponse(.success(session)))
                } catch: { error, send in
                    print("Start session failed → \(error)")
                    await send(.sessionStartResponse(.failure(error)))
                }

            case let .sessionStartResponse(.success(session)):
                return .run { send in
                    print("Sending startRound request → sessionId: \(session.id)")
                    let round = try await sessionAPIClient.startRound(session.id)
                    print("Start round success → \(round.id)")
                    await send(.roundStartResponse(.success(round)))
                } catch: { error, send in
                    print("Start round failed → \(error)")
                    await send(.roundStartResponse(.failure(error)))
                }

                
            case let .roundStartResponse(.success(round)):
                state.timer.currentRoundId = round.id

                let work = (round.workTime ?? 25) * 60         // ← RoundResult の値を使う
                let shortBreak = (round.breakTime ?? 5) * 60   // ← 同上
                let longBreak = state.timer.longBreakDuration
                let roundsPerSession = state.timer.roundsPerSession

                return .merge(
                    .send(.timer(.updateSettings(
                        task: work,
                        shortBreak: shortBreak,
                        longBreak: longBreak,
                        roundsPerSession: roundsPerSession
                    ))),
                    .send(.timer(.start))
                )

                
            case .timer(.phaseCompleted):
                if state.timer.phase == .shortBreak || state.timer.phase == .longBreak {
                    state.evalModal = EvalModalFeature.State(
                        score: 0.5,
                        round: state.timer.round
                    )
                }
                return .none
                
            case .evalModal(.submitEval(let score)):
                state.evalModal = nil
                guard let roundId = state.timer.currentRoundId else {
                    print("⚠️ roundId is nil")
                    return .none
                }

                print("📨 評価送信中: roundId=\(roundId), score=\(score)")

                return .run { send in
                    let result = try await sessionAPIClient.completeRound(roundId, Int(score * 100))
                    await send(.completeRoundResponse(.success(result)))
                } catch: { error, send in
                    print("❌ completeRound エラー: \(error)")
                    await send(.completeRoundResponse(.failure(error)))
                }

            case .startNextRound:
                guard let sessionId = state.timer.sessionId else {
                    print("⚠️ sessionId is nil")
                    return .none
                }
                return .run { send in
                    let round = try await sessionAPIClient.startRound(sessionId)
                    print("🔄 Next round started → \(round.id)")
                    await send(.roundStartResponse(.success(round)))
                } catch: { error, send in
                    print("❌ Failed to start next round: \(error)")
                    await send(.roundStartResponse(.failure(error)))
                }

            case let .completeRoundResponse(.success(round)):
                print("✅ completeRound 成功: \(round)")
                return .send(.timer(.start))
                
            case .dismissEvalModal:
                state.evalModal = nil
                return .none
                
            default:
                return .none
            }
        }
    }
}
