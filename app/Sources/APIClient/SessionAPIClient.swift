//
//  SessionAPIClient.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/05/21.
//

import Foundation
import Dependencies

struct SessionAPIClient {
    var startSession: () async throws -> SessionResult
    var completeSession: (_ sessionId: UUID) async throws -> SessionResult
    var startRound: (_ sessionId: UUID) async throws -> RoundResult
    var completeRound: (_ roundId: UUID, _ focusScore: Int) async throws -> RoundResult
    var getSession: (_ sessionId: UUID) async throws -> SessionResult
}

extension SessionAPIClient {
    static let live = SessionAPIClient(
        startSession: {
            var request = URLRequest(url: URL(string: "http://localhost:8080/api/v1/sessions")!)
            request.httpMethod = "POST"
            request.setValue("Bearer dev-token", forHTTPHeaderField: "Authorization")
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")

            let (data, _) = try await URLSession.shared.data(for: request)
            print("Start session → \(String(data: data, encoding: .utf8) ?? "")")
            return try APIFormatters.jsonDecoder.decode(SessionResult.self, from: data)
        },

        completeSession: { sessionId in
            var request = URLRequest(url: URL(string: "http://localhost:8080/api/v1/sessions/\(sessionId)/complete")!)
            request.httpMethod = "PATCH"
            request.setValue("Bearer dev-token", forHTTPHeaderField: "Authorization")
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")

            let (data, _) = try await URLSession.shared.data(for: request)
            print("Complete session → \(String(data: data, encoding: .utf8) ?? "")")
            return try APIFormatters.jsonDecoder.decode(SessionResult.self, from: data)
        },

        startRound: { sessionId in
            var request = URLRequest(url: URL(string: "http://localhost:8080/api/v1/sessions/\(sessionId)/rounds")!)
            request.httpMethod = "POST"
            request.setValue("Bearer dev-token", forHTTPHeaderField: "Authorization")
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")

            let (data, _) = try await URLSession.shared.data(for: request)
            print("Start round → \(String(data: data, encoding: .utf8) ?? "")")
            return try APIFormatters.jsonDecoder.decode(RoundResult.self, from: data)
        },

        completeRound: { roundId, focusScore in
            var request = URLRequest(url: URL(string: "http://localhost:8080/api/v1/rounds/\(roundId)/complete")!)
            request.httpMethod = "PATCH"
            request.setValue("Bearer dev-token", forHTTPHeaderField: "Authorization")
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")
            request.httpBody = try JSONEncoder().encode(["focus_score": focusScore])

            let (data, _) = try await URLSession.shared.data(for: request)
            print("Complete round → \(String(data: data, encoding: .utf8) ?? "")")
            return try APIFormatters.jsonDecoder.decode(RoundResult.self, from: data)
        },

        getSession: { sessionId in
            var request = URLRequest(url: URL(string: "http://localhost:8080/api/v1/sessions/\(sessionId)")!)
            request.httpMethod = "GET"
            request.setValue("Bearer dev-token", forHTTPHeaderField: "Authorization")
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")

            let (data, _) = try await URLSession.shared.data(for: request)
            return try APIFormatters.jsonDecoder.decode(SessionResult.self, from: data)
        }
    )
}

extension DependencyValues {
    var sessionAPIClient: SessionAPIClient {
        get { self[SessionAPIClientKey.self] }
        set { self[SessionAPIClientKey.self] = newValue }
    }

    private enum SessionAPIClientKey: DependencyKey {
        static let liveValue = SessionAPIClient.live
    }
}
