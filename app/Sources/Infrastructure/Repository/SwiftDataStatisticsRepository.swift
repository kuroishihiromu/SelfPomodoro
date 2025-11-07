//
//  SwiftDataStatisticsRepository.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/07.
//

import Foundation
import SwiftData

@MainActor
final class SwiftDataStatisticsRepository: StatisticsRepository {
    private let context: ModelContext
    private let calendar = Calendar.current

    init(context: ModelContext) {
        self.context = context
    }

    func fetchFocusTrend(forWeekStarting startDate: Date, userIdentifier: String) async throws -> [ConcentrationData] {
        let weekStart = calendar.startOfDay(for: startDate)
        guard let weekEnd = calendar.date(byAdding: .day, value: 6, to: weekStart) else {
            return []
        }
        let historyStart = calendar.date(byAdding: .day, value: -30, to: weekStart) ?? weekStart

        var descriptor = FetchDescriptor<RoundRecordModel>(
            predicate: #Predicate {
                $0.userIdentifier == userIdentifier &&
                $0.createdAt >= historyStart && $0.createdAt <= weekEnd && $0.focusScore != nil
            }
        )
        descriptor.sortBy = [SortDescriptor(\.createdAt, order: .forward)]
        let models = try context.fetch(descriptor)

        let normalized = aggregateDailyScores(models: models)
        let chartData = ChartDataProcessor.calculateMovingAverage(from: normalized)
        let filtered = chartData.filter {
            $0.date >= weekStart && $0.date <= weekEnd
        }

        // Ensure there are 7 entries (fill missing days)
        var resultDict = Dictionary(uniqueKeysWithValues: filtered.map { ($0.date, $0) })
        let weekDates = (0..<7).compactMap { calendar.date(byAdding: .day, value: $0, to: weekStart) }
        return weekDates.map { date in
            resultDict[date] ?? ConcentrationData(date: date, score: 0, movingAverage: 0, stdDev: 0)
        }
    }

    func fetchHeatMapData(forMonth month: Date, userIdentifier: String) async throws -> [FocusData] {
        let monthStart = calendar.date(from: calendar.dateComponents([.year, .month], from: month))!
        guard let monthEnd = calendar.date(byAdding: DateComponents(month: 1, day: -1), to: monthStart) else {
            return []
        }

        var descriptor = FetchDescriptor<RoundRecordModel>(
            predicate: #Predicate {
                $0.userIdentifier == userIdentifier &&
                $0.createdAt >= monthStart && $0.createdAt <= monthEnd && $0.focusScore != nil
            }
        )
        descriptor.sortBy = [SortDescriptor(\.createdAt, order: .forward)]
        let models = try context.fetch(descriptor)

        return models.compactMap { model in
            guard let score = model.focusScore else { return nil }
            return FocusData(
                date: model.createdAt,
                hour: calendar.component(.hour, from: model.createdAt),
                focus_score: score
            )
        }
    }

    private func aggregateDailyScores(models: [RoundRecordModel]) -> [FocusTrendResult] {
        let grouped = Dictionary(grouping: models) { calendar.startOfDay(for: $0.createdAt) }
        return grouped.keys.sorted().map { day in
            let scores = grouped[day]?.compactMap(\.focusScore) ?? []
            let average = scores.isEmpty ? 0 : Double(scores.reduce(0, +)) / Double(scores.count)
            return FocusTrendResult(date: day, focusScore: average)
        }
    }
}
