//
//  StatisticsAPIClient.swift
//  SelfPomodoro
//
//  Created by し on 2025/06/15.
//

import Foundation
import Dependencies
import Amplify


enum StatisticsAPIError: Error, Equatable {
    case networkError
    case decodingError
    case unknown
}

struct StatisticsAPIClient {
    var fetchConcentrationData: () async throws -> [ConcentrationData]
}

extension StatisticsAPIClient {
    static let live = StatisticsAPIClient(
        fetchConcentrationData: {
            let request = RESTRequest(
                apiName: "selfpomodoro",
                path: "/dev/api/v1/statistics/focus-trend"
            )

            do {
                let data = try await Amplify.API.get(request: request)

                let focusResults = try AppDecoder.default.decode(FocusTrendResponse.self, from: data)

                // 移動平均と標準偏差の計算（直近7日）
                let windowSize = 7
                var concentrationDataList: [ConcentrationData] = []

                for (index, result) in focusResults.items.enumerated() {
                    let start = max(0, index - windowSize + 1)
                    let window = focusResults.items[start...index].map { $0.focusScore }

                    let average = window.reduce(0, +) / Double(window.count)
                    let variance = window.map { pow($0 - average, 2) }.reduce(0, +) / Double(window.count)
                    let stdDev = sqrt(variance)

                    let data = ConcentrationData(
                        date: result.date,
                        score: result.focusScore,
                        movingAverage: average,
                        stdDev: stdDev
                    )
                    concentrationDataList.append(data)
                }

                return concentrationDataList
            } catch let decodingError as DecodingError {
                throw StatisticsAPIError.decodingError
            } catch {
                throw StatisticsAPIError.networkError
            }
        }
    )
}

private enum StatisticsAPIClientKey: DependencyKey {
    static let liveValue = StatisticsAPIClient.live
}

extension DependencyValues {
    var statisticsAPIClient: StatisticsAPIClient {
        get { self[StatisticsAPIClientKey.self] }
        set { self[StatisticsAPIClientKey.self] = newValue }
    }
}
