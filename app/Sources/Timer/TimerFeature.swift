//
//  TimerFeature.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/04/09.
//

import ComposableArchitecture
import Foundation

@Reducer
struct TimerFeature {

    enum Phase: Equatable {
        case task
        case shortBreak
        case longBreak
    }

    struct State: Equatable {
        var currentRoundId: UUID?
        var sessionId: UUID?
        var currentSeconds: Int = 0
        var totalSeconds: Int

        var isRunning: Bool = false
        var round: Int = 1

        var phase: Phase = .task

        var taskDuration: Int
        var shortBreakDuration: Int
        var longBreakDuration: Int

        var roundsPerSession: Int

        var startTime: ContinuousClock.Instant? = nil
        var lastTaskDuration: Int = 0

        var currentPhaseDuration: Int {
            switch phase {
            case .task:
                return taskDuration
            case .shortBreak:
                return shortBreakDuration
            case .longBreak:
                return longBreakDuration
            }
        }

    }

    enum Action: Equatable {
        case start
        case stop
        case tick(Int)
        case phaseCompleted(completedPhase: Phase)
        case updateSettings(task: Int, shortBreak: Int, longBreak: Int, roundsPerSession: Int)
        case saveTimerState
        case restoreTimerState
        case clearPersistedState
    }

    enum CancelID { case timer }

    func reduce(into state: inout State, action: Action) -> Effect<Action> {
        switch action {

        case .start:
            state.isRunning = true
            state.totalSeconds = state.currentPhaseDuration
            let correctedStart = ContinuousClock().now.advanced(by: .seconds(-state.currentSeconds))
            state.startTime = correctedStart

            // 永続化
            let persistenceData = TimerPersistenceData(
                startTime: Date(timeIntervalSinceNow: -Double(state.currentSeconds)),
                taskDuration: state.taskDuration,
                shortBreakDuration: state.shortBreakDuration,
                longBreakDuration: state.longBreakDuration,
                roundsPerSession: state.roundsPerSession,
                phase: phaseToString(state.phase),
                round: state.round,
                isRunning: true,
                currentSeconds: state.currentSeconds,
                sessionId: state.sessionId,
                currentRoundId: state.currentRoundId
            )
            TimerPersistence.save(persistenceData)
            return .run { [start = correctedStart] send in
                var lastElapsed = -1
                while !Task.isCancelled {
                    let now = ContinuousClock().now
                    let realElapsed = start.duration(to: now).components.seconds

                    let accelerationFactor = 20.0  // リリース時は通常速度

                    let acceleratedElapsed = Int(Double(realElapsed) * accelerationFactor)

                    if acceleratedElapsed != lastElapsed {
                        await send(.tick(acceleratedElapsed))
                        lastElapsed = acceleratedElapsed
                    }

                    try? await Task.sleep(nanoseconds: 100_000_000)  // 0.1秒ごとにチェック（=リアルタイム）
                }
            }
            .cancellable(id: CancelID.timer)

        case .stop:
            state.isRunning = false
            TimerPersistence.clear()
            return .cancel(id: CancelID.timer)

        case let .tick(elapsed):
            guard elapsed != state.currentSeconds else {
                return .none
            }
            state.currentSeconds = elapsed
            if elapsed >= state.totalSeconds {
                return .send(.phaseCompleted(completedPhase: state.phase))
            }
            return .none

        case let .phaseCompleted(completedPhase):
            state.isRunning = false
            let elapsed = state.currentSeconds
            if state.phase == .task {
                state.lastTaskDuration = min(elapsed, state.currentPhaseDuration)
            }
            state.currentSeconds = 0
            switch state.phase {
            case .task:
                // セッションの最後のタスクだった場合は longBreak
                if state.round >= state.roundsPerSession {
                    state.phase = .longBreak
                } else {
                    state.phase = .shortBreak
                }
            case .shortBreak:
                state.phase = .task
                state.round += 1

            case .longBreak:
                state.phase = .task
                state.round = 1

            }

            state.totalSeconds = state.currentPhaseDuration
            TimerPersistence.clear()
            return .send(.stop)

        case let .updateSettings(task, short, long, rps):
            print(
                "🛠️ Timer updateSettings: from task=\(state.taskDuration), short=\(state.shortBreakDuration), long=\(state.longBreakDuration), rps=\(state.roundsPerSession) -> to task=\(task), short=\(short), long=\(long), rps=\(rps)"
            )
            state.taskDuration = task
            state.shortBreakDuration = short
            state.longBreakDuration = long
            state.roundsPerSession = rps
            state.totalSeconds = state.currentPhaseDuration
            state.currentSeconds = 0
            return .none

        case .saveTimerState:
            guard state.isRunning else { return .none }
            let persistenceData = TimerPersistenceData(
                // 現在の経過秒に合わせた論理開始時刻を保存
                startTime: Date(timeIntervalSinceNow: -Double(state.currentSeconds)),
                taskDuration: state.taskDuration,
                shortBreakDuration: state.shortBreakDuration,
                longBreakDuration: state.longBreakDuration,
                roundsPerSession: state.roundsPerSession,
                phase: phaseToString(state.phase),
                round: state.round,
                isRunning: state.isRunning,
                currentSeconds: state.currentSeconds,
                sessionId: state.sessionId,
                currentRoundId: state.currentRoundId
            )
            TimerPersistence.save(persistenceData)
            return .none

        case .restoreTimerState:
            guard let persistedData = TimerPersistence.load() else {
                return .none
            }

            state.taskDuration = persistedData.taskDuration
            state.shortBreakDuration = persistedData.shortBreakDuration
            state.longBreakDuration = persistedData.longBreakDuration
            state.roundsPerSession = persistedData.roundsPerSession
            state.phase = stringToPhase(persistedData.phase)
            state.round = persistedData.round
            state.sessionId = persistedData.sessionId
            state.currentRoundId = persistedData.currentRoundId
            state.totalSeconds = state.currentPhaseDuration

            print(
                "🔄 Timer state restored - sessionId: \(persistedData.sessionId?.uuidString ?? "nil"), currentRoundId: \(persistedData.currentRoundId?.uuidString ?? "nil")"
            )

            if persistedData.isRunning {
                // アプリ非稼働中の経過は実時間のみ反映（加速は適用しない）
                let elapsed = max(0, Int(Date().timeIntervalSince(persistedData.startTime)))
                let adjustedElapsed: Int = elapsed

                if adjustedElapsed < state.totalSeconds {
                    state.currentSeconds = adjustedElapsed
                    state.isRunning = true
                    let correctedStart = ContinuousClock().now.advanced(
                        by: .seconds(-adjustedElapsed))
                    state.startTime = correctedStart
                    return .send(.start)
                } else {
                    TimerPersistence.clear()
                    return .send(.phaseCompleted(completedPhase: state.phase))
                }
            }

            return .none

        case .clearPersistedState:
            TimerPersistence.clear()
            return .none
        }
    }

    private func phaseToString(_ phase: Phase) -> String {
        switch phase {
        case .task: return "task"
        case .shortBreak: return "shortBreak"
        case .longBreak: return "longBreak"
        }
    }

    private func stringToPhase(_ string: String) -> Phase {
        switch string {
        case "task": return .task
        case "shortBreak": return .shortBreak
        case "longBreak": return .longBreak
        default: return .task
        }
    }
}
