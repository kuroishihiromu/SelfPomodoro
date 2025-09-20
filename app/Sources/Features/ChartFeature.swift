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
                let targetDate = state.currentWeekStart
                return .run { send in
                    print("📈 ChartFeature: fetchFocusTrend start")
                    do {
                        let data = try await apiClient.fetchConcentrationData(targetDate)
                        print("📈 ChartFeature: fetchFocusTrend success count=\(data.count)")
                        await send(.dataLoaded(data))
                    } catch {
                        print("📉 ChartFeature: fetchFocusTrend failed: \(error)")
                    }
                }

            case let .dataLoaded(data):
                print("🧮 ChartFeature: dataLoaded count=\(data.count)")
                state.data = data
                return .none
                
            case .previousWeek:
                state.currentWeekStart = Calendar.current.date(byAdding: .day, value: -7, to: state.currentWeekStart)!
                return .send(.fetchFocusTrend)

            case .nextWeek:
                state.currentWeekStart = Calendar.current.date(byAdding: .day, value: 7, to: state.currentWeekStart)!
                return .send(.fetchFocusTrend)
            }
        }
    }
}
