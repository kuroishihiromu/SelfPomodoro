//
//  StatisticsAPIClient.swift
//  SelfPomodoro
//
//  Created by し on 2025/06/15.
//

import Foundation
import Dependencies

struct StatisticsAPIClient {
    var fetchFocusTrend: () async throws -> [FocusTrendResult]
}

extension StatisticsAPIClient {
    static let live = StatisticsAPIClient(
        fetchFocusTrend: {
            guard let url = Bundle.main.url(forResource: "focus-trend", withExtension: "json") else {
                throw URLError(.badURL)
            }
            let data = try Data(contentsOf: url)
            return try AppDecoder.default.decode([FocusTrendResult].self, from: data)
        }
    )
}

extension DependencyValues {
    var statisticsAPIClient: StatisticsAPIClient {
        get { self[StatisticsAPIClientKey.self] }
        set { self[StatisticsAPIClientKey.self] = newValue }
    }

    private enum StatisticsAPIClientKey: DependencyKey {
        static let liveValue = StatisticsAPIClient.live
    }
}
