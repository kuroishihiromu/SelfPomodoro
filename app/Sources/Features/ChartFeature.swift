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
    @Dependency(\.statisticsAPIClient) var apiClient

    struct State: Equatable {
        var data: [ConcentrationData] = []
        var currentWeekStart: Date = .startOfCurrentWeek()

        var weekDates: [Date] {
            (0..<8).compactMap { offset in
                Calendar.current.date(byAdding: .day, value: offset, to: currentWeekStart)
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

    var body: some ReducerOf<Self> {
        Reduce { state, action in
            switch action {

            case .fetchFocusTrend:
                return .run { send in
                    do {
                        let focusTrendResults = try await apiClient.fetchFocusTrend()
                        let normalizeToDayStart = ChartDataProcessor.normalizeDates(focusTrendResults)
                        let scoredConcentrationData = ChartDataProcessor.calculateMovingAverage(from: normalizeToDayStart)
                        await send(.dataLoaded(scoredConcentrationData))
                    } catch {
                        print("データ読み込み失敗: \(error)")
                    }
                }

            case let .dataLoaded(data):
                state.data = data
                return .none

            case .previousWeek:
                state.currentWeekStart = Calendar.current.date(byAdding: .day, value: -7, to: state.currentWeekStart)!
                return .none

            case .nextWeek:
                state.currentWeekStart = Calendar.current.date(byAdding: .day, value: +7, to: state.currentWeekStart)!
                return .none
            }
        }
    }
}
