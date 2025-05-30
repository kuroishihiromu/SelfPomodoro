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
        case phaseCompleted
        case updateSettings(task: Int, shortBreak: Int, longBreak: Int, roundsPerSession: Int)
    }

    enum CancelID { case timer }

    func reduce(into state: inout State, action: Action) -> Effect<Action> {
        switch action {

        case .start:
            state.isRunning = true
            state.totalSeconds = state.currentPhaseDuration
            let correctedStart = ContinuousClock().now.advanced(by: .seconds(-state.currentSeconds))
            state.startTime = correctedStart
            return .run { [start = correctedStart] send in
                var lastElapsed = -1
                while !Task.isCancelled {
                    let now = ContinuousClock().now
                    let elapsed = Int(start.duration(to: now).components.seconds)
                    if elapsed != lastElapsed {
                        await send(.tick(elapsed))
                        lastElapsed = elapsed
                    }
                    try? await Task.sleep(nanoseconds: 100_000_000)
                }
            }
            .cancellable(id: CancelID.timer)

        case .stop:
            state.isRunning = false
            return .cancel(id: CancelID.timer)

        case let .tick(elapsed):
            guard elapsed != state.currentSeconds else {
                return .none
            }
            state.currentSeconds = elapsed
            if elapsed >= state.totalSeconds {
                return .send(.phaseCompleted)
            }
            return .none

        case .phaseCompleted:
            state.isRunning = false
            state.currentSeconds = 0

            switch state.phase {
            case .task:
                // セッションの最後のタスクだった場合は longBreak
                if state.round % state.roundsPerSession == 0 {
                    state.phase = .longBreak
                } else {
                    state.phase = .shortBreak
                }
            case .shortBreak, .longBreak:
                state.phase = .task
                state.round += 1
            }

            state.totalSeconds = state.currentPhaseDuration
            return .send(.stop)

        case let .updateSettings(task, short, long, rps):
            state.taskDuration = task
            state.shortBreakDuration = short
            state.longBreakDuration = long
            state.roundsPerSession = rps
            state.totalSeconds = state.currentPhaseDuration
            state.currentSeconds = 0
            return .none
        }
    }
}
