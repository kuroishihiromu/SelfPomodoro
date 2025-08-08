//
//  ChartFeature.swift
//  SelfPomodoro
//
//  Created by し on 2025/06/29.
//

import ComposableArchitecture
import Foundation

@Reducer
struct ChartFeature {

    @ObservableState
    struct State: Equatable {
        var data: [ConcentrationData] = []
        var currentWeekStart: Date = .startOfCurrentWeek()

        var weekDates: [Date] {
            (0..<7).compactMap {
                Calendar.current.date(byAdding: .day, value: $0, to: currentWeekStart)
            }
        }

        var currentWeekData: [ConcentrationData] {
            let dict = Dictionary(uniqueKeysWithValues: data.map { ($0.date, $0) })
            return weekDates.map { date in
                dict[date] ?? ConcentrationData(date: date, score: 0, movingAverage: 0, stdDev: 0)
            }
        }
    }

    enum Action: Equatable {
        case fetchFocusTrend
        case dataLoaded([ConcentrationData])
        case previousWeek
        case nextWeek
    }

    @Dependency(\.statisticsAPIClient) var apiClient

    var body: some ReducerOf<Self> {
        Reduce { state, action in
            switch action {
            case .fetchFocusTrend:
                return .run { send in
                    do {
                        let data = try await apiClient.fetchConcentrationData()
                        await send(.dataLoaded(data))
                    } catch {
                        print("データ取得失敗: \(error)")
                    }
                }

            case let .dataLoaded(data):
                state.data = data
                state.currentWeekStart = .startOfCurrentWeek()
                return .none
                
            case .previousWeek:
                state.currentWeekStart = Calendar.current.date(byAdding: .day, value: -7, to: state.currentWeekStart)!
                return .none

            case .nextWeek:
                state.currentWeekStart = Calendar.current.date(byAdding: .day, value: 7, to: state.currentWeekStart)!
                return .none
            }
        }
    }
}
