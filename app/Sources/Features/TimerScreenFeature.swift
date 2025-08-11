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
    
    @ObservableState
    struct State: Equatable {
        var timer: TimerFeature.State
        var evalModal: EvalModalFeature.State?
        var sessionId: UUID?
        var roundConfigModalIsPresented: Bool = false
        var sessionCompleteModal: Bool = false
        var isFirstSession: Bool = true
        var userConfig: UserConfigResult = .init(id: UUID(), roundWorkTime: 25*60, roundBreakTime: 5*60, sessionRounds: 5, sessionBreakTime: 15*60)
    }

    enum Action {
        case timer(TimerFeature.Action)
        case evalModal(EvalModalFeature.Action)

        case StartRoundButtonTapped

        case startNextRound
        
        case sessionStartResponse(Result<SessionResult, Error>)
        case roundStartResponse(Result<RoundResult, Error>)
        case completeRoundResponse(Result<RoundResult, Error>)

        case showEvalModal
        case dismissEvalModal
        case sessionCompleteModalTapped
        case onAppear
        case refreshUserConfig
        case userConfigResponse(Result<UserConfigResult, Error>)
        case toggleConfigModal(Bool)
        case toggleSessionCompleteModal(Bool)
        
    }

    @Dependency(\.sessionAPIClient) var sessionAPIClient
    @Dependency(\.userConfigAPIClient) var userConfigAPIClient
    
    var body: some ReducerOf<Self> {
        Scope(state: \.timer, action: \.timer) { TimerFeature() }
        .ifLet(\.evalModal, action: \.evalModal) { EvalModalFeature() }

        Reduce { state, action in
            switch action {

            case .StartRoundButtonTapped:
                
                if !state.isFirstSession {
                    return .send(.startNextRound)
                } else {
                    return .run { send in
                        let session = try await sessionAPIClient.startSession()
                        await send(.sessionStartResponse(.success(session)))
                    } catch: { error, send in
                        await send(.sessionStartResponse(.failure(error)))
                    }
                }

            case let .sessionStartResponse(.success(session)):
                state.sessionId = session.id
                state.timer.sessionId = session.id
                state.isFirstSession = false
                return .send(.startNextRound)

            case let .roundStartResponse(.success(round)):
                state.timer.currentRoundId = round.id

                return .send(.timer(.start))

            case .timer(.phaseCompleted):
                if state.timer.phase == .shortBreak || state.timer.phase == .longBreak {
                    state.evalModal = EvalModalFeature.State(
                        score: 0.5,
                        round: state.timer.round
                    )
                    
                } else if state.timer.phase == .task {
                    if state.userConfig.sessionRounds < state.timer.round {
                        state.sessionCompleteModal = true
                        state.timer.round = 1
                    } else {
                        state.roundConfigModalIsPresented = true
                    }
                }
                return .none

            case .evalModal(.submitEval(let score)):
                state.evalModal = nil
                guard let roundId = state.timer.currentRoundId else {
                    return .none
                }

                return .run { send in
                    let result = try await sessionAPIClient.completeRound(roundId, Int(score * 100))
                    await send(.completeRoundResponse(.success(result)))
                } catch: { error, send in
                    await send(.completeRoundResponse(.failure(error)))
                }

            case .startNextRound:
                guard let sessionId = state.timer.sessionId else {
                    return .none
                }
                return .run { send in
                    let round = try await sessionAPIClient.startRound(sessionId)
                    await send(.roundStartResponse(.success(round)))
                } catch: { error, send in
                    await send(.roundStartResponse(.failure(error)))
                }

            case let .completeRoundResponse(.success(round)):
                return .send(.timer(.start))

            case .completeRoundResponse(.failure(let error)):
                return .none

            case .dismissEvalModal:
                state.evalModal = nil
                return .none

            case .onAppear:
                guard state.isFirstSession else {
                    return .none
                }
                
                return .run { send in
                    let config = try await userConfigAPIClient.getUserConfig()
                    await send(.userConfigResponse(.success(config)))
                } catch: { error, send in
                    await send(.userConfigResponse(.failure(error)))
                }

            case let .userConfigResponse(.success(config)):
                state.userConfig = config
                state.roundConfigModalIsPresented = true
                
                return .send(.timer(.updateSettings(
                    task: config.roundWorkTime * 60,
                    shortBreak: config.roundBreakTime * 60,
                    longBreak: config.sessionBreakTime * 60,
                    roundsPerSession: config.sessionRounds
                )))

            case let .toggleConfigModal(show):
                state.roundConfigModalIsPresented = show
                return .none
                
            case let .toggleSessionCompleteModal(show):
                state.sessionCompleteModal = show
                if !show {
                    return .send(.toggleConfigModal(true))
                }
                return .none

            case .sessionCompleteModalTapped:
                state.isFirstSession = true
                return .none
                
            default:
                return .none
            }
        }
    }
}
