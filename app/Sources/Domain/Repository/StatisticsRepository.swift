//
//  StatisticsRepository.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/07.
//

import Foundation

protocol StatisticsRepository {
    func fetchFocusTrend(forWeekStarting startDate: Date, userIdentifier: String) async throws -> [ConcentrationData]
    func fetchHeatMapData(forMonth month: Date, userIdentifier: String) async throws -> [FocusData]
}
