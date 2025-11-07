//
//  OptimizationAPIClient.swift
//  SelfPomodoro
//
//  Created by Codex on 2025/02/15.
//

import Dependencies
import Foundation

struct OptimizationAPIClient {
    var sendRoundData: (_ userId: UUID, _ records: [OptimizationRoundPayload]) async throws -> OptimizationRoundResponse?
    var sendSessionData: (_ userId: UUID, _ records: [OptimizationSessionPayload]) async throws -> OptimizationSessionResponse?
}

struct OptimizationRoundPayload: Codable, Equatable {
    let time: String
    let work_time: Double
    let break_time: Double
    let focus_score: Int
}

struct OptimizationSessionPayload: Codable, Equatable {
    let time: String
    let round_count: Int
    let total_work_time: Double
    let break_time: Double
    let avg_focus_score: Double
}

struct OptimizationRoundResponse: Decodable {
    let workTime: Double
    let breakTime: Double
}

struct OptimizationSessionResponse: Decodable {
    let totalWorkTime: Double
    let breakTime: Double
    let roundCount: Int
}

extension OptimizationAPIClient {
    static let live = OptimizationAPIClient(
        sendRoundData: { userId, records in
            guard !records.isEmpty else { return nil }
            let url = URL(string: "https://6cacxhgjoz2bejfydr3vogjpla0bwvkb.lambda-url.ap-northeast-1.on.aws/round/v4/\(userId.uuidString)")!
            var request = URLRequest(url: url)
            request.httpMethod = "POST"
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")
            request.httpBody = try JSONEncoder().encode(records)

            let (data, response) = try await URLSession.shared.data(for: request)
            guard let httpResponse = response as? HTTPURLResponse, 200..<300 ~= httpResponse.statusCode else {
                throw OptimizationAPIError.invalidResponse
            }
            print("📮 Round optimization sent, bytes=\(data.count)")
            if let raw = String(data: data, encoding: .utf8) {
                print("📦 Round optimization raw response: \(raw)")
            }
            let decoder = JSONDecoder()
            decoder.keyDecodingStrategy = .convertFromSnakeCase
            do {
                return try decoder.decode(OptimizationRoundResponse.self, from: data)
            } catch {
                print("⚠️ Round optimization decode failed: \(error)")
                return nil
            }
        },
        sendSessionData: { userId, records in
            guard !records.isEmpty else { return nil }
            let url = URL(string: "https://6cacxhgjoz2bejfydr3vogjpla0bwvkb.lambda-url.ap-northeast-1.on.aws/session/v3/\(userId.uuidString)")!
            var request = URLRequest(url: url)
            request.httpMethod = "POST"
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")
            request.httpBody = try JSONEncoder().encode(records)

            let (data, response) = try await URLSession.shared.data(for: request)
            guard let httpResponse = response as? HTTPURLResponse, 200..<300 ~= httpResponse.statusCode else {
                throw OptimizationAPIError.invalidResponse
            }
            print("📮 Session optimization sent, bytes=\(data.count)")
            if let raw = String(data: data, encoding: .utf8) {
                print("📦 Session optimization raw response: \(raw)")
            }
            let decoder = JSONDecoder()
            decoder.keyDecodingStrategy = .convertFromSnakeCase
            do {
                return try decoder.decode(OptimizationSessionResponse.self, from: data)
            } catch {
                print("⚠️ Session optimization decode failed: \(error)")
                return nil
            }
        }
    )
}

enum OptimizationAPIError: Error {
    case invalidResponse
}

extension DependencyValues {
    var optimizationAPIClient: OptimizationAPIClient {
        get { self[OptimizationAPIClientKey.self] }
        set { self[OptimizationAPIClientKey.self] = newValue }
    }

    private enum OptimizationAPIClientKey: DependencyKey {
        static let liveValue = OptimizationAPIClient.live
    }
}
