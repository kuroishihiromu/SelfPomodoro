//
//  ChartDataProcessor.swift
//  SelfPomodoro
//
//  Created by し on 2025/06/29.
//

import Foundation

enum ChartDataProcessor {
    static func normalizeDates(_ results: [FocusTrendResult]) -> [FocusTrendResult] {
        let calendar = Calendar.current
        return results.map { result in
            let normalizedDate = calendar.startOfDay(for: result.date)
            return FocusTrendResult(date: normalizedDate, focusScore: result.focusScore)
        }
    }

    static func calculateMovingAverage(from results: [FocusTrendResult], windowSize: Int = 7) -> [ConcentrationData] {
        let sorted = results.sorted { $0.date < $1.date }
        let scores = sorted.map { $0.focusScore }

        return sorted.enumerated().map { index, result in
            let start = max(0, index - windowSize + 1)
            let window = scores[start...index]
            let average = window.reduce(0, +) / Double(window.count)
            let variance = window.map { pow($0 - average, 2) }.reduce(0, +) / Double(window.count)
            let stdDev = sqrt(variance)

            return ConcentrationData(
                date: result.date,
                score: scores[index],
                movingAverage: average,
                stdDev: stdDev
            )
        }
    }
}
