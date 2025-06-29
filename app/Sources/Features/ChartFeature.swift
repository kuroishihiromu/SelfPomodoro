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
        var currentWeekStart: Date = Self.startOfCurrentWeek()

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

        static func startOfCurrentWeek() -> Date {
            let calendar = Calendar.current
            let today = calendar.startOfDay(for: Date())
            let weekday = calendar.component(.weekday, from: today)
            let diff = (weekday + 5) % 7
            return calendar.date(byAdding: .day, value: -diff, to: today)!
        }
    }

    enum Action: Equatable {
        case fetchData
        case dataLoaded([ConcentrationData])
        case previousWeek
        case nextWeek
    }

    var body: some ReducerOf<Self> {
        Reduce { state, action in
            switch action {

            case .fetchData:
                return .run { send in
                    do {
                        let raw = try await apiClient.fetchFocusTrend()
                        let normalized = ChartDataProcessor.normalizeDates(raw)
                        let processed = ChartDataProcessor.calculateMovingAverage(from: normalized)
                        await send(.dataLoaded(processed))
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

    struct ConcentrationData: Identifiable, Equatable {
        var id: Date { date }
        let date: Date
        let score: Double
        let movingAverage: Double
        let stdDev: Double
    }
}
