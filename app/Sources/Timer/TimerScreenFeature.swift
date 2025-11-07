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
        // ユーザー設定は分単位で保持。タイマー適用時に秒へ変換。
        var userConfig: UserConfig = .default()
        var hasPersistedTimer: Bool = false
        var didRestoreFromPersistence: Bool = false
        var currentSessionRounds: [RoundRecord] = []
        var toast: ToastState = ToastState(message: "")
    }

    enum Action {
        case timer(TimerFeature.Action)
        case evalModal(EvalModalFeature.Action)

        case StartRoundButtonTapped

        case startNextRound

        case completeRoundResponse(Result<RoundRecord, Error>)
        case completeSessionResponse(Result<SessionRecord, Error>)

        case showEvalModal
        case dismissEvalModal
        case sessionCompleteModalTapped
        case onAppear
        case refreshUserConfig
        case userConfigResponse(Result<UserConfig, Error>)
        case toggleConfigModal(Bool)
        case toggleSessionCompleteModal(Bool)
        case restoreTimerIfNeeded
        case toastDismissed
        case showToast(String)
        
    }

    @Dependency(\.roundRecordRepository) var roundRecordRepository
    @Dependency(\.sessionRecordRepository) var sessionRecordRepository
    @Dependency(\.userConfigRepository) var userConfigRepository
    @Dependency(\.optimizationAPIClient) var optimizationAPIClient
    @Dependency(\.userIdentifier) var resolveUserIdentifier
    
    var body: some ReducerOf<Self> {
        Scope(state: \.timer, action: \.timer) { TimerFeature() }
        .ifLet(\.evalModal, action: \.evalModal) { EvalModalFeature() }

        Reduce { state, action in
            switch action {

            case .StartRoundButtonTapped:
                if state.isFirstSession || state.sessionId == nil {
                    let newSessionId = UUID()
                    state.sessionId = newSessionId
                    state.timer.sessionId = newSessionId
                    state.currentSessionRounds = []
                    state.isFirstSession = false
                    print("🎆 Session started locally: \(newSessionId)")
                }
                state.roundConfigModalIsPresented = false
                return .send(.startNextRound)

            case .startNextRound:
                let newRoundId = UUID()
                state.timer.currentRoundId = newRoundId
                print("🎯 Starting round id=\(newRoundId)")
                return .merge(
                    .send(.timer(.saveTimerState)),
                    .send(.timer(.start))
                )

            case let .timer(.phaseCompleted(completedPhase)):
                if completedPhase == .task {
                    // task完了時は評価モーダルを表示
                    state.evalModal = EvalModalFeature.State(
                        score: 0.5,
                        round: state.timer.round
                    )
                } else if completedPhase == .shortBreak {
                    // shortBreak完了 → 次のラウンドのtaskへ
                    state.roundConfigModalIsPresented = true
                    return .send(.refreshUserConfig)
                } else if completedPhase == .longBreak {
                    // longBreak完了 → 新しいセッション開始可能
                    state.isFirstSession = true
                    state.roundConfigModalIsPresented = true
                    return .send(.refreshUserConfig)
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
                guard let sessionId = state.timer.sessionId else {
                    print("⚠️ Session ID missing, cannot save round")
                    return .none
                }

                let roundId = state.timer.currentRoundId ?? UUID()
                state.timer.currentRoundId = roundId

                let identifier = resolveUserIdentifier()
                let workMinutes = Double(state.timer.lastTaskDuration) / 60.0
                let breakMinutes: Double = state.timer.round >= state.userConfig.sessionRounds
                    ? state.userConfig.sessionBreakMinutes
                    : state.userConfig.roundBreakMinutes
                let focusScore = Int(score * 100)
                let record = RoundRecord(
                    id: roundId,
                    userIdentifier: identifier,
                    sessionIdentifier: sessionId,
                    workMinutes: workMinutes,
                    breakMinutes: breakMinutes,
                    focusScore: focusScore,
                    isAborted: false,
                    createdAt: Date(),
                    updatedAt: Date()
                )

                return .run { send in
                    try await roundRecordRepository.save(record)
                    await send(.completeRoundResponse(.success(record)))
                } catch: { error, send in
                    await send(.completeRoundResponse(.failure(error)))
                }

            case let .completeRoundResponse(.success(round)):
                state.currentSessionRounds.append(round)
                let roundOptimizationEffect = sendRoundOptimization(for: round.userIdentifier)

                if state.timer.round >= state.userConfig.sessionRounds {
                    state.sessionCompleteModal = true
                    guard let sessionId = state.timer.sessionId else {
                        return .none
                    }

                    let identifier = resolveUserIdentifier()
                    let rounds = state.currentSessionRounds
                    let sessionRecord = SessionRecord(
                        id: UUID(),
                        userIdentifier: identifier,
                        sessionIdentifier: sessionId,
                        roundCount: rounds.count,
                        totalWorkMinutes: rounds.reduce(0) { $0 + $1.workMinutes },
                        breakMinutes: state.userConfig.sessionBreakMinutes,
                        averageFocusScore: averageFocus(from: rounds),
                        isAborted: false,
                        createdAt: Date(),
                        updatedAt: Date()
                    )

                    state.isFirstSession = true
                    state.sessionId = nil
                    state.timer.sessionId = nil
                    state.currentSessionRounds = []
                    let sessionOptimizationEffect = sendSessionOptimization(for: identifier)

                    return .merge(
                        .run { send in
                        try await sessionRecordRepository.save(sessionRecord)
                        await send(.completeSessionResponse(.success(sessionRecord)))
                    } catch: { error, send in
                        await send(.completeSessionResponse(.failure(error)))
                    },
                        sessionOptimizationEffect,
                        roundOptimizationEffect
                    )
                } else {
                    return .merge(
                        .send(.timer(.start)),
                        roundOptimizationEffect
                    )
                }

            case .completeRoundResponse(.failure(let error)):
                print("⚠️ Round record save failed: \(error)")
                return .none

            case .dismissEvalModal:
                state.evalModal = nil
                return .none

            case .onAppear:
                // 既に実行中なら何もしない（タブ復帰時の二重スタート防止）
                if state.timer.isRunning {
                    return .none
                }

                let persistedExists = TimerPersistence.load() != nil
                state.hasPersistedTimer = persistedExists

                // 一度だけ復元する
                if persistedExists && !state.didRestoreFromPersistence {
                    return .send(.restoreTimerIfNeeded)
                }

                if state.isFirstSession {
                    return .send(.refreshUserConfig)
                }

                return .none

            case .refreshUserConfig:
                let identifier = resolveUserIdentifier()
                return .run { send in
                    if let config = try await userConfigRepository.fetchLatest(for: identifier) {
                        await send(.userConfigResponse(.success(config)))
                    } else {
                        let defaults = UserConfig.default(for: identifier)
                        try await userConfigRepository.upsertLatest(defaults)
                        await send(.userConfigResponse(.success(defaults)))
                    }
                } catch: { error, send in
                    print("⚠️ Fetch user config failed: \(error)")
                    await send(.userConfigResponse(.failure(error)))
                }

            case let .userConfigResponse(.success(config)):
                state.userConfig = config
                print("🧭 Applying UserConfig to timer (minutes): work=\(config.roundWorkMinutes), break=\(config.roundBreakMinutes), rounds=\(config.sessionRounds), longBreak=\(config.sessionBreakMinutes)")
                // 永続化データがない場合のみモーダルを表示
                if !state.hasPersistedTimer {
                    state.roundConfigModalIsPresented = true
                }
                
                // サーバーは分単位を返すため、タイマー内部の秒に変換して適用
                return .send(.timer(.updateSettings(
                    task: Int(config.roundWorkMinutes * 60),
                    shortBreak: Int(config.roundBreakMinutes * 60),
                    longBreak: Int(config.sessionBreakMinutes * 60),
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
                return .none

            case let .completeSessionResponse(.success(record)):
                print("✅ SessionRecord saved id=\(record.id)")
                return .send(.refreshUserConfig)

            case let .completeSessionResponse(.failure(error)):
                print("⚠️ Session record save failed: \(error)")
                return .none

            case .sessionCompleteModalTapped:
                // セッション完了モーダルの"休憩を開始"ボタンが押された
                // モーダルを閉じてlongBreakを開始
                state.sessionCompleteModal = false
                return .send(.timer(.start))
                
            case .restoreTimerIfNeeded:
                // 永続化データから sessionId を復元
                if let persistedData = TimerPersistence.load() {
                    print("🔄 Restoring session state: sessionId=\(persistedData.sessionId?.uuidString ?? "nil"), currentRoundId=\(persistedData.currentRoundId?.uuidString ?? "nil")")
                    state.sessionId = persistedData.sessionId
                    state.hasPersistedTimer = true
                    state.didRestoreFromPersistence = true
                    state.roundConfigModalIsPresented = false
                    
                    // セッションが復元された場合は初回セッションではない
                    if persistedData.sessionId != nil {
                        state.isFirstSession = false
                    }
                }
                
                return .send(.timer(.restoreTimerState))
                
            case .toastDismissed:
                state.toast.isVisible = false
                return .none
                
            case let .showToast(message):
                state.toast.message = message
                state.toast.isVisible = true
                return .run { send in
                    try? await Task.sleep(nanoseconds: 2_000_000_000)
                    await send(.toastDismissed)
                }
                
            default:
                return .none
            }
        }
    }

    private func averageFocus(from rounds: [RoundRecord]) -> Double {
        let scores = rounds.compactMap(\.focusScore)
        guard !scores.isEmpty else { return 0 }
        return Double(scores.reduce(0, +)) / Double(scores.count)
    }

    private func sendRoundOptimization(for userIdentifier: String) -> Effect<Action> {
        guard let uuid = UUID(uuidString: userIdentifier) else {
            print("⚠️ Invalid userIdentifier for optimization: \(userIdentifier)")
            return .none
        }

        return .run { send in
            let records = try await roundRecordRepository.fetchRecent(for: userIdentifier, limit: nil)
            let payloads = records.compactMap { record -> OptimizationRoundPayload? in
                guard let focus = record.focusScore else { return nil }
                return OptimizationRoundPayload(
                    time: Self.isoFormatter.string(from: record.createdAt),
                    work_time: record.workMinutes,
                    break_time: record.breakMinutes,
                    focus_score: focus
                )
            }
            guard !payloads.isEmpty else { return }
            do {
                if let response = try await optimizationAPIClient.sendRoundData(uuid, payloads) {
                    let work = response.workTime
                    let rest = response.breakTime
                    print("📬 Round optimization response: work=\(work), break=\(rest)")
                    let entry = UserConfigRoundHistoryEntry(
                        id: UUID(),
                        userIdentifier: userIdentifier,
                        recommendedWorkMinutes: work,
                        recommendedBreakMinutes: rest,
                        createdAt: Date()
                    )
                    try await userConfigRepository.addRoundHistory(entry)

                    var latest = try await userConfigRepository.fetchLatest(for: userIdentifier) ?? UserConfig.default(for: userIdentifier)
                    latest.roundWorkMinutes = work
                    latest.roundBreakMinutes = rest
                    latest.updatedAt = Date()
                    try await userConfigRepository.upsertLatest(latest)
                    
                    await send(.showToast("ラウンドの最適化が完了しました"))
                }
            } catch {
                print("⚠️ Round optimization failed: \(error)")
            }
        }
    }

    private func sendSessionOptimization(for userIdentifier: String) -> Effect<Action> {
        guard let uuid = UUID(uuidString: userIdentifier) else {
            print("⚠️ Invalid userIdentifier for session optimization: \(userIdentifier)")
            return .none
        }

        return .run { send in
            let records = try await sessionRecordRepository.fetchRecent(for: userIdentifier, limit: nil)
            let payloads = records.map { record in
                OptimizationSessionPayload(
                    time: Self.isoFormatter.string(from: record.createdAt),
                    round_count: record.roundCount,
                    total_work_time: record.totalWorkMinutes,
                    break_time: record.breakMinutes,
                    avg_focus_score: record.averageFocusScore
                )
            }
            guard !payloads.isEmpty else { return }
            do {
                if let response = try await optimizationAPIClient.sendSessionData(uuid, payloads) {
                    let rounds = response.roundCount
                    let rest = response.breakTime
                    print("📬 Session optimization response: rounds=\(rounds), break=\(rest)")
                    let entry = UserConfigSessionHistoryEntry(
                        id: UUID(),
                        userIdentifier: userIdentifier,
                        recommendedSessionRounds: rounds,
                        recommendedSessionBreakMinutes: rest,
                        createdAt: Date()
                    )
                    try await userConfigRepository.addSessionHistory(entry)

                    var latest = try await userConfigRepository.fetchLatest(for: userIdentifier) ?? UserConfig.default(for: userIdentifier)
                    latest.sessionRounds = rounds
                    latest.sessionBreakMinutes = rest
                    latest.updatedAt = Date()
                    try await userConfigRepository.upsertLatest(latest)
                    
                    await send(.showToast("セッションの最適化が完了しました"))
                }
            } catch {
                print("⚠️ Session optimization failed: \(error)")
            }
        }
    }

    private static let isoFormatter: ISO8601DateFormatter = {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        return formatter
    }()
}
