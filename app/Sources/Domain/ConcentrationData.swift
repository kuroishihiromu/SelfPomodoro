//
//  ConcentrationData.swift
//  SelfPomodoro
//
//  Created by し on 2025/07/03.
//

import Foundation

struct ConcentrationData: Identifiable, Equatable {
    var id: Date { date }
    let date: Date
    let score: Double
    let movingAverage: Double
    let stdDev: Double
}
