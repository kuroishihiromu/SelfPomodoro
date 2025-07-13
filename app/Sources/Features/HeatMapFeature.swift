//
//  HeatMapFeature.swift
//  SelfPomodoro
//
//  Created by し on 2025/07/13.
//

import ComposableArchitecture
import Foundation

struct HeatMapFeature: Reducer {
    struct State: Equatable {
        var currentMonth: Date = Date()
        var focusData: [FocusData] = []
    }

    enum Action: Equatable {
        case fetchHeatMap
        case heatMapDataLoaded([FocusData])
        case previousMonth
        case nextMonth
    }

    func reduce(into state: inout State, action: Action) -> Effect<Action> {
        switch action {
        case .fetchHeatMap:
            let data = DummyDataLoader.loadFocusData()
            return .send(.heatMapDataLoaded(data))

        case let .heatMapDataLoaded(data):
            state.focusData = data
            return .none

        case .previousMonth:
            state.currentMonth = Calendar.current.date(byAdding: .month, value: -1, to: state.currentMonth) ?? state.currentMonth
            return .send(.fetchHeatMap)

        case .nextMonth:
            state.currentMonth = Calendar.current.date(byAdding: .month, value: 1, to: state.currentMonth) ?? state.currentMonth
            return .send(.fetchHeatMap)
        }
    }
}
