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
        var hasPersistedTimer: Bool = false
    }

    enum Action {
        case timer(TimerFeature.Action)
        case evalModal(EvalModalFeature.Action)

        case StartRoundButtonTapped

        case startNextRound
        
        case sessionStartResponse(Result<SessionResult, Error>)
        case roundStartResponse(Result<RoundResult, Error>)
        case completeRoundResponse(Result<RoundResult, Error>)
        case completeSessionResponse(Result<SessionResult, Error>)

        case showEvalModal
        case dismissEvalModal
        case sessionCompleteModalTapped
        case onAppear
        case refreshUserConfig
        case userConfigResponse(Result<UserConfigResult, Error>)
        case toggleConfigModal(Bool)
        case toggleSessionCompleteModal(Bool)
        case restoreTimerIfNeeded
        
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
                print("🎆 Session started: \(session.id)")
                state.sessionId = session.id
                state.timer.sessionId = session.id
                state.isFirstSession = false
                return .send(.startNextRound)

            case let .sessionStartResponse(.failure(error)):
                // 409 Conflict（セッションが既に存在）の場合は既存セッションとして次のラウンドへ進む
                let message = String(describing: error)
                if message.contains("409") {
                    state.isFirstSession = false
                    return .send(.startNextRound)
                }
                return .none

            case let .roundStartResponse(.success(round)):
                print("🎯 roundStartResponse: setting currentRoundId to \(round.id)")
                state.timer.currentRoundId = round.id
                
                // ラウンド開始時に状態を永続化
                return .merge(
                    .send(.timer(.saveTimerState)),
                    .send(.timer(.start))
                )

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
                        if let sessionId = state.timer.sessionId {
                            // セッション完了をサーバに通知
                            return .run { send in
                                let result = try await sessionAPIClient.completeSession(sessionId)
                                await send(.completeSessionResponse(.success(result)))
                            } catch: { error, send in
                                await send(.completeSessionResponse(.failure(error)))
                            }
                        }
                    } else {
                        state.roundConfigModalIsPresented = true
                        // ラウンドが切り替わるタイミングでユーザー設定を再取得
                        return .send(.refreshUserConfig)
                    }
                }
                return .none

            case let .completeSessionResponse(.success(session)):
                // サーバ側でセッションは完了済み。次回は新規セッションを開始できるように状態をリセット
                state.timer.sessionId = nil
                state.sessionId = nil
                state.isFirstSession = true
                return .none

            case .completeSessionResponse(.failure):
                // エラー時は UI を妨げない（後で再試行できるようにする）
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
                // 永続化されたタイマーデータをチェック
                state.hasPersistedTimer = TimerPersistence.load() != nil
                
                if state.hasPersistedTimer {
                    // 永続化データがある場合は復元
                    return .send(.restoreTimerIfNeeded)
                } else if state.isFirstSession {
                    // 永続化データがなく、初回セッションの場合は設定を取得
                    return .send(.refreshUserConfig)
                }
                
                return .none

            case .refreshUserConfig:
                return .run { send in
                    print("🔄 Fetching user config...")
                    let config = try await userConfigAPIClient.getUserConfig()
                    await send(.userConfigResponse(.success(config)))
                } catch: { error, send in
                    print("⚠️ Fetch user config failed: \(error)")
                    await send(.userConfigResponse(.failure(error)))
                }

            case let .userConfigResponse(.success(config)):
                state.userConfig = config
                print("🧭 Applying UserConfig to timer: work=\(config.roundWorkTime), break=\(config.roundBreakTime), rounds=\(config.sessionRounds), longBreak=\(config.sessionBreakTime)")
                // 永続化データがない場合のみモーダルを表示
                if !state.hasPersistedTimer {
                    state.roundConfigModalIsPresented = true
                }
                
                // サーバーは秒単位を返す前提。×60 せずにそのまま適用。
                return .send(.timer(.updateSettings(
                    task: config.roundWorkTime,
                    shortBreak: config.roundBreakTime,
                    longBreak: config.sessionBreakTime,
                    roundsPerSession: config.sessionRounds
                )))

            case let .userConfigResponse(.failure(error)):
                print("❗️ userConfigResponse failure: \(error)")
                return .none

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
                
            case .restoreTimerIfNeeded:
                // 永続化データから sessionId を復元
                if let persistedData = TimerPersistence.load() {
                    print("🔄 Restoring session state: sessionId=\(persistedData.sessionId?.uuidString ?? "nil"), currentRoundId=\(persistedData.currentRoundId?.uuidString ?? "nil")")
                    state.sessionId = persistedData.sessionId
                    state.hasPersistedTimer = true
                    state.roundConfigModalIsPresented = false
                    
                    // セッションが復元された場合は初回セッションではない
                    if persistedData.sessionId != nil {
                        state.isFirstSession = false
                    }
                }
                
                return .send(.timer(.restoreTimerState))
                
            default:
                return .none
            }
        }
    }
}
