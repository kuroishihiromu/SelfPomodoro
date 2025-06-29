//
//  StatisticsFeature.swift
//  SelfPomodoro
//
//  Created by し on 2025/06/15.
//

import Foundation

struct ConcentrationData: Identifiable {
    var id: Date { date }
    let date: Date
    let score: Double
    let movingAverage: Double
    let stdDev: Double
}

func normalizeDates(_ results: [FocusTrendResult]) -> [FocusTrendResult] {
    let calendar = Calendar.current
    return results.map { result in
        let normalizedDate = calendar.startOfDay(for: result.date)
        return FocusTrendResult(
            date: normalizedDate,
            focusScore: result.focusScore
        )
    }
}

func calculateMovingAverage(from results: [FocusTrendResult], windowSize: Int = 7) -> [ConcentrationData] {
    let sorted = results.sorted { $0.date < $1.date }
    let scores = sorted.map { Double($0.focusScore) }

    var output: [ConcentrationData] = []

    for i in 0..<sorted.count {
        let startIdx = max(0, i - windowSize + 1)
        let window = scores[startIdx...i]
        let average = window.reduce(0, +) / Double(window.count)
        let variance = window.map { pow($0 - average, 2) }.reduce(0, +) / Double(window.count)
        let stdDev = sqrt(variance)
        let data = ConcentrationData(
            date: sorted[i].date,
            score: scores[i],
            movingAverage: average,
            stdDev: stdDev
        )
        output.append(data)
    }

    return output
}
