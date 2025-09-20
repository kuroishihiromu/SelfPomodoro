//
//  ConfigAPIClient.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/07/03.
//

import Foundation
import Dependencies
import Amplify
import AWSPluginsCore

enum userconfigAPIError: Error, Equatable {
    case networkError
    case decodingError
    case unknown
}

struct UserConfigAPIClient {
    var getUserConfig: () async throws -> UserConfigResult
}

extension UserConfigAPIClient {
    static let live = UserConfigAPIClient(
        getUserConfig: {
            let idToken: String
            do {
                let session = try await Amplify.Auth.fetchAuthSession()
                guard let provider = session as? AuthCognitoTokensProvider else {
                    throw userconfigAPIError.unknown
                }
                let tokens = try provider.getCognitoTokens().get()
                idToken = tokens.idToken
                print("👤 Auth session (fetchUserconfig) isSignedIn=\(session.isSignedIn)")
            } catch {
                print("👤 Auth session (fetchUserconfig) fetch failed: \(error)")
                throw error
            }
            print("➡️ GET /dev/api/v1/optimization-preferences")
            let request = RESTRequest(
                apiName: "selfpomodoro",
                path: "/dev/api/v1/optimization-preferences",
                headers: ["Authorization" : idToken]
            )
            
            do {
                let data = try await Amplify.API.get(request: request)
                let userconfig = try APIFormatters.jsonDecoder.decode(UserConfigResult.self, from: data)
                return userconfig
            } catch {
                print("❌ fetchUserconfig failed: \(error)")
                throw error
            }
        }
    )
}

extension DependencyValues {
    var userConfigAPIClient: UserConfigAPIClient {
        get { self[UserConfigAPIClientKey.self] }
        set { self[UserConfigAPIClientKey.self] = newValue }
    }

    private enum UserConfigAPIClientKey: DependencyKey {
        static let liveValue = UserConfigAPIClient.live
    }
}
