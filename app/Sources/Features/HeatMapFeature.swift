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

struct FocusData: Equatable, Identifiable, Decodable {
    let id = UUID()
    let date: Date
    let hour: Int
    let focus_score: Int
}

struct FocusDataWrapper: Decodable {
    let items: [FocusData]
}

enum DummyDataLoader {
    static func loadFocusData() -> [FocusData] {
        guard let url = Bundle.main.url(forResource: "heatmap_data", withExtension: "json"),
              let data = try? Data(contentsOf: url) else {
            return []
        }

        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .formatted(DateFormatter.yyyyMMdd)

        if let decoded = try? decoder.decode(FocusDataWrapper.self, from: data) {
            return decoded.items
        }

        return []
    }
}

private extension DateFormatter {
    static let yyyyMMdd: DateFormatter = {
        let formatter = DateFormatter()
        formatter.dateFormat = "yyyy-MM-dd"
        formatter.timeZone = .current
        return formatter
    }()
}
