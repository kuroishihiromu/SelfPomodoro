//
//  StatisticsFeature.swift
//  SelfPomodoro
//
//  Created by し on 2025/06/15.
//

import ComposableArchitecture
import Foundation

struct StatisticsFeature: Reducer {
    struct State: Equatable {
        var chart = ChartFeature.State()
    }

    @CasePathable
    enum Action: Equatable {
        case chart(ChartFeature.Action)
    }

    var body: some ReducerOf<Self> {
        Scope(state: \ .chart, action: \ .chart) {
            ChartFeature()
        }
    }
}
