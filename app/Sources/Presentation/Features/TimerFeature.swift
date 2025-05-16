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

    struct State: Equatable {
        var pomodoro: PomodoroTimer
        var startTime: ContinuousClock.Instant? = nil
    }
    
    struct TimerEnvironment {
        var startUseCase: StartTimerUseCase
        var stopUseCase: StopTimerUseCase
        var tickUseCase: TickTimerUseCase
        var advancePhaseUseCase: AdvancePhaseUseCase
        var updateSettingsUseCase: UpdateSettingsUseCase
    }

    enum Action: Equatable {
        case start
        case stop
        case tick(Int)
        case phaseCompleted
        case updateSettings(task: Int, shortBreak: Int, longBreak: Int, roundsPerSession: Int)
    }

    enum CancelID { case timer }

    func reduce(into state: inout State, action: Action, env: TimerEnvironment) -> Effect<Action> {
        switch action {

        case .start:
            state.pomodoro = env.startUseCase.execute(timer: state.pomodoro)
            let correctedStart = ContinuousClock().now.advanced(by: .seconds(-state.pomodoro.currentSeconds))
            state.startTime = correctedStart
            return .run { [start = correctedStart, timer = state.pomodoro] send in
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
            state.pomodoro = env.stopUseCase.execute(timer: state.pomodoro)
            return .cancel(id: CancelID.timer)

        case let .tick(elapsed):
            guard elapsed != state.pomodoro.currentSeconds else {
                return .none
            }
            state.pomodoro.currentSeconds = elapsed
            if elapsed >= state.pomodoro.totalSeconds {
                return .send(.phaseCompleted)
            }
            return .none

        case .phaseCompleted:
            state.pomodoro = env.advancePhaseUseCase.execute(timer: state.pomodoro)
            state.pomodoro.currentSeconds = 0
            return .send(.stop)

        case let .updateSettings(task, short, long, rps):
            state.pomodoro = env.updateSettingsUseCase.execute(
                task: task,
                short: short,
                long: long,
                rounds: rps
            )
            return .none
        }
    }
}
