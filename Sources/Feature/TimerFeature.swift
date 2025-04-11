//
//  TimerFeature.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/04/09.
//

import Foundation
import ComposableArchitecture

struct TimerFeature: Reducer {
    struct State: Equatable {
        var currentSeconds: Int
        var totalSeconds: Int
        var isRunning: Bool = false
    }

    enum Action: Equatable {
        case start
        case stop
        case tick
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
                state.isRunning = false
                return .cancel(id: CancelID.timer)
            }
            state.currentSeconds += 1
            return .none
        }
    }
}
