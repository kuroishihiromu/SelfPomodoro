//
//  SessionAPIClient.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/05/21.
//

import Foundation
import Dependencies
import Amplify
import AWSPluginsCore

enum SessionAPIError: Error, Equatable {
    case networkError
    case decodingError
    case unknown
}

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
            let idToken: String
            do {
                let session = try await Amplify.Auth.fetchAuthSession()
                guard let provider = session as? AuthCognitoTokensProvider else {
                    throw SessionAPIError.unknown
                }
                let tokens = try provider.getCognitoTokens().get()
                idToken = tokens.idToken
                print("👤 Auth session (startSession) isSignedIn=\(session.isSignedIn)")
            } catch {
                print("👤 Auth session (startSession) fetch failed: \(error)")
                throw error
            }
            print("➡️ POST /dev/api/v1/sessions")
            let request = RESTRequest(
                apiName: "selfpomodoro",
                path: "/dev/api/v1/sessions",
                headers: ["Authorization" : idToken]
            )

            do {
                let data = try await Amplify.API.post(request: request)
                print("📦 startSession bytes=\(data.count)")
                let result = try APIFormatters.jsonDecoder.decode(SessionResult.self, from: data)
                print("✅ startSession id=\(result.id)")
                return result
            } catch {
                print("❌ startSession failed: \(error)")
                throw error
            }
        },

        completeSession: { sessionId in
            let idToken: String
            do {
                let session = try await Amplify.Auth.fetchAuthSession()
                guard let provider = session as? AuthCognitoTokensProvider else {
                    throw SessionAPIError.unknown
                }
                let tokens = try provider.getCognitoTokens().get()
                idToken = tokens.idToken
                print("👤 Auth session (completeSession) isSignedIn=\(session.isSignedIn)")
            } catch {
                print("👤 Auth session (completeSession) fetch failed: \(error)")
                throw error
            }
            print("➡️ PATCH /dev/api/v1/sessions/\(sessionId)/complete")
            let request = RESTRequest(
                apiName: "selfpomodoro",
                path: "/dev/api/v1/sessions/\(sessionId)/complete",
                headers: ["Authorization" : idToken]
            )

            do {
                let data = try await Amplify.API.patch(request: request)
                print("📦 completeSession bytes=\(data.count)")
                let result = try APIFormatters.jsonDecoder.decode(SessionResult.self, from: data)
                print("✅ completeSession id=\(result.id)")
                return result
            } catch {
                print("❌ completeSession failed: \(error)")
                throw error
            }
        },
        
        startRound: { sessionId in
            let idToken: String
            do {
                let session = try await Amplify.Auth.fetchAuthSession()
                guard let provider = session as? AuthCognitoTokensProvider else {
                    throw SessionAPIError.unknown
                }
                let tokens = try provider.getCognitoTokens().get()
                idToken = tokens.idToken
                print("👤 Auth session (startRound) isSignedIn=\(session.isSignedIn)")
            } catch {
                print("👤 Auth session (startRound) fetch failed: \(error)")
                throw error
            }
            let sessionIdLower = sessionId.uuidString.lowercased()
            print("➡️ POST /dev/api/v1/sessions/\(sessionIdLower)/rounds")
            let request = RESTRequest(
                apiName: "selfpomodoro",
                path: "/dev/api/v1/sessions/\(sessionIdLower)/rounds",
                headers: ["Authorization" : idToken]
            )

            do {
                let data = try await Amplify.API.post(request: request)
                print("📦 startRound bytes=\(data.count)")
                let result = try APIFormatters.jsonDecoder.decode(RoundResult.self, from: data)
                print("✅ startRound id=\(result.id) order=\(result.roundOrder)")
                return result
            } catch {
                print("❌ startRound failed: \(error)")
                throw error
            }
        },
        
        completeRound: { roundId, focusScore in
            let idToken: String
            do {
                let session = try await Amplify.Auth.fetchAuthSession()
                guard let provider = session as? AuthCognitoTokensProvider else {
                    throw SessionAPIError.unknown
                }
                let tokens = try provider.getCognitoTokens().get()
                idToken = tokens.idToken
                print("👤 Auth session (completeRound) isSignedIn=\(session.isSignedIn)")
            } catch {
                print("👤 Auth session (completeRound) fetch failed: \(error)")
                throw error
            }
            print("➡️ PATCH /dev/api/v1/rounds/\(roundId)/complete")
            let body = try JSONEncoder().encode(["focus_score": focusScore])
            let request = RESTRequest(
                apiName: "selfpomodoro",
                path: "/dev/api/v1/rounds/\(roundId)/complete",
                headers: ["Authorization" : idToken],
                body: body
            )

            do {
                let data = try await Amplify.API.patch(request: request)
                print("📦 completeRound bytes=\(data.count)")
                let result = try APIFormatters.jsonDecoderWithISOEasyVersion.decode(RoundResult.self, from: data)
                print("✅ completeRound id=\(result.id) focus=\(result.focusScore ?? -1)")
                print("complete Round body:::: \(result)")
                return result
            } catch {
                print("❌ completeRound failed: \(error)")
                throw error
            }
        },

        getSession: { sessionId in
            let idToken: String
            do {
                let session = try await Amplify.Auth.fetchAuthSession()
                guard let provider = session as? AuthCognitoTokensProvider else {
                    throw SessionAPIError.unknown
                }
                let tokens = try provider.getCognitoTokens().get()
                idToken = tokens.idToken
                print("👤 Auth session (getSession) isSignedIn=\(session.isSignedIn)")
            } catch {
                print("👤 Auth session (getSession) fetch failed: \(error)")
                throw error
            }
            print("➡️ GET /dev/api/v1/sessions/\(sessionId)")
            let request = RESTRequest(
                apiName: "selfpomodoro",
                path: "/dev/api/v1/sessions/\(sessionId)",
                headers: ["Authorization" : idToken]
            )

            do {
                let data = try await Amplify.API.get(request: request)
                print("📦 getSession bytes=\(data.count)")
                let result = try APIFormatters.jsonDecoder.decode(SessionResult.self, from: data)
                print("✅ getSession id=\(result.id)")
                return result
            } catch {
                print("❌ getSession failed: \(error)")
                throw error
            }
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
