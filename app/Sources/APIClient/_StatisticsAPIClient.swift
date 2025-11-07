// //
// //  StatisticsAPIClient.swift
// //  SelfPomodoro
// //
// //  Created by し on 2025/06/15.
// //

// import AWSPluginsCore
// import Amplify
// import Dependencies
// import Foundation

// enum StatisticsAPIError: Error, Equatable {
//     case networkError
//     case decodingError
//     case unknown
// }

// struct StatisticsAPIClient {
//     var fetchConcentrationData: (_ date: Date) async throws -> [ConcentrationData]
// }

// extension StatisticsAPIClient {
//     static let live = StatisticsAPIClient(
//         fetchConcentrationData: { date in
//             let idToken: String
//             do {
//                 let session = try await Amplify.Auth.fetchAuthSession()
//                 guard let provider = session as? AuthCognitoTokensProvider else {
//                     throw StatisticsAPIError.unknown
//                 }
//                 let tokens = try provider.getCognitoTokens().get()
//                 idToken = tokens.idToken
//                 print(
//                     "👤 Auth session (statistics.fetchConcentrationData) isSignedIn=\(session.isSignedIn)"
//                 )
//             } catch {
//                 print("👤 Auth session (statistics.fetchConcentrationData) fetch failed: \(error)")
//                 throw error
//             }
//             // yyyyMMdd 形式で日付パラメータを付与
//             let comps = Calendar.current.dateComponents([.year, .month, .day], from: date)
//             let yyyymmdd = String(
//                 format: "%04d%02d%02d", comps.year ?? 0, comps.month ?? 0, comps.day ?? 0)
//             print("➡️ GET /dev/api/v1/statistics/focus-trend/\(yyyymmdd)")
//             let request = RESTRequest(
//                 apiName: "selfpomodoro",
//                 path: "/dev/api/v1/statistics/focus-trend/\(yyyymmdd)",
//                 headers: ["Authorization": "Bearer \(idToken)"]
//             )

//             do {
//                 let data = try await Amplify.API.get(request: request)
//                 print("📦 statistics.fetchConcentrationData bytes=\(data.count)")

//                 // サーバーが配列 or { items: [] } どちらでも対応
//                 let resultsArray: [FocusTrendResult]
//                 if let arr = try? AppDecoder.default.decode([FocusTrendResult].self, from: data) {
//                     resultsArray = arr
//                     print("✅ statistics.fetchConcentrationData items=\(arr.count) (top-level array)")
//                 } else {
//                     let wrapper = try AppDecoder.default.decode(FocusTrendResponse.self, from: data)
//                     resultsArray = wrapper.items
//                     print("✅ statistics.fetchConcentrationData items=\(resultsArray.count) (wrapped)")
//                 }

//                 // 移動平均と標準偏差の計算（直近7日）
//                 let windowSize = 7
//                 var concentrationDataList: [ConcentrationData] = []

//                 for (index, result) in resultsArray.enumerated() {
//                     let start = max(0, index - windowSize + 1)
//                     let window = resultsArray[start...index].map { $0.focusScore }

//                     let average = window.reduce(0, +) / Double(window.count)
//                     let variance =
//                         window.map { pow($0 - average, 2) }.reduce(0, +) / Double(window.count)
//                     let stdDev = sqrt(variance)

//                     let data = ConcentrationData(
//                         date: result.date,
//                         score: result.focusScore,
//                         movingAverage: average,
//                         stdDev: stdDev
//                     )
//                     concentrationDataList.append(data)
//                 }

//                 return concentrationDataList
//             } catch let decodingError as DecodingError {
//                 print("❌ statistics.fetchConcentrationData decode failed: \(decodingError)")
//                 throw StatisticsAPIError.decodingError
//             } catch {
//                 print("❌ statistics.fetchConcentrationData failed: \(error)")
//                 throw StatisticsAPIError.networkError
//             }
//         }
//     )
// }

// private enum StatisticsAPIClientKey: DependencyKey {
//     static let liveValue = StatisticsAPIClient.live
// }

// extension DependencyValues {
//     var statisticsAPIClient: StatisticsAPIClient {
//         get { self[StatisticsAPIClientKey.self] }
//         set { self[StatisticsAPIClientKey.self] = newValue }
//     }
// }
