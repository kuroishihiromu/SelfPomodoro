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
        case rest
    }

    struct State: Equatable {
        var currentSeconds: Int = 0
        var totalSeconds: Int

        var isRunning: Bool = false
        var round: Int = 1

        var phase: Phase = .task

        var taskDuration: Int
        var restDuration: Int

        // 現在のフェーズに対応する時間を再取得する
        var currentPhaseDuration: Int {
            phase == .task ? taskDuration : restDuration
        }
    }

    enum Action: Equatable {
        case start
        case stop
        case tick
        case phaseCompleted
        case updateDurations(task: Int, rest: Int)  // ← 今後API対応するならコレ
    }

    enum CancelID { case timer }

    func reduce(into state: inout State, action: Action) -> Effect<Action> {
        switch action {

        case .start:
            state.isRunning = true
            return .run { send in
                while true {
                    try await Task.sleep(nanoseconds: 1_000_000_000)
                    await send(.tick)
                }
            }
            .cancellable(id: CancelID.timer)

        case .stop:
            state.isRunning = false
            return .cancel(id: CancelID.timer)

        case .tick:
            guard state.currentSeconds < state.totalSeconds else {
                return .send(.phaseCompleted)
            }
            state.currentSeconds += 1
            return .none

        case .phaseCompleted:
            state.isRunning = false
            state.currentSeconds = 0

            // フェーズ切り替えと round 増加
            switch state.phase {
            case .task:
                state.phase = .rest
            case .rest:
                state.phase = .task
                state.round += 1
            }

            // 次のフェーズの時間を設定
            state.totalSeconds = state.currentPhaseDuration
            return .send(.stop)

        case let .updateDurations(task, rest):
            state.taskDuration = task
            state.restDuration = rest

            // 今のフェーズに合わせて totalSeconds を更新
            state.totalSeconds = state.currentPhaseDuration
            state.currentSeconds = 0
            return .none
        }
    }
}
