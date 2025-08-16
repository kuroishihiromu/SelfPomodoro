//
//  SessionAPIClient.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/05/21.
//

import Foundation
import Dependencies
import Amplify

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
            @Dependency(\.tokenProvider) var tokenProvider
            
            var headers = [
                "Content-Type": "application/json"
            ]
            
            if let token = tokenProvider.getAuthToken() {
                headers["Authorization"] = "Bearer \(token)"
            }
            
            let request = RESTRequest(
                apiName: "selfpomodoro",
                path: "/dev/api/v1/sessions",
                headers: headers
            )

            let data = try await Amplify.API.post(request: request)
            return try APIFormatters.jsonDecoder.decode(SessionResult.self, from: data)
        },

        completeSession: { sessionId in
            @Dependency(\.tokenProvider) var tokenProvider
            
            var headers = [
                "Content-Type": "application/json"
            ]
            
            if let token = tokenProvider.getAuthToken() {
                headers["Authorization"] = "Bearer \(token)"
            }
            
            let request = RESTRequest(
                apiName: "selfpomodoro",
                path: "/dev/api/v1/sessions/\(sessionId)/complete",
                headers: headers
            )

            let data = try await Amplify.API.patch(request: request)
            return try APIFormatters.jsonDecoder.decode(SessionResult.self, from: data)
        },
        
        startRound: { sessionId in
            @Dependency(\.tokenProvider) var tokenProvider
            
            var headers = [
                "Content-Type": "application/json"
            ]
            
            if let token = tokenProvider.getAuthToken() {
                headers["Authorization"] = "Bearer \(token)"
            }
            
            let sessionIdLower = sessionId.uuidString.lowercased()
            let request = RESTRequest(
                apiName: "selfpomodoro",
                path: "/dev/api/v1/sessions/\(sessionIdLower)/rounds",
                headers: headers
            )

            let data = try await Amplify.API.post(request: request)
            return try APIFormatters.jsonDecoder.decode(RoundResult.self, from: data)
        },
        
        completeRound: { roundId, focusScore in
            @Dependency(\.tokenProvider) var tokenProvider
            
            var headers = [
                "Content-Type": "application/json"
            ]
            
            if let token = tokenProvider.getAuthToken() {
                headers["Authorization"] = "Bearer \(token)"
            }
            
            let body = try JSONEncoder().encode(["focus_score": focusScore])
            let request = RESTRequest(
                apiName: "selfpomodoro",
                path: "/dev/api/v1/rounds/\(roundId)/complete",
                headers: headers,
                body: body
            )

            let data = try await Amplify.API.patch(request: request)
            return try APIFormatters.jsonDecoderWithISOEasyVersion.decode(RoundResult.self, from: data)
        },

        getSession: { sessionId in
            @Dependency(\.tokenProvider) var tokenProvider
            
            var headers = [
                "Content-Type": "application/json"
            ]
            
            if let token = tokenProvider.getAuthToken() {
                headers["Authorization"] = "Bearer \(token)"
            }
            
            let request = RESTRequest(
                apiName: "selfpomodoro",
                path: "/dev/api/v1/sessions/\(sessionId)",
                headers: headers
            )

            let data = try await Amplify.API.get(request: request)
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
