//
//  AuthAPIClient.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/06/16.
//

import Foundation
import Amplify
import AWSCognitoAuthPlugin
import AWSPluginsCore
import Dependencies

struct AuthTokens: Equatable {
    let idToken: String
    let accessToken: String
    let refreshToken: String
}

struct AuthAPIClient {
    var signIn: (_ username: String, _ password: String) async throws -> AuthTokens
    var signOut: () async throws -> Void
}

extension DependencyValues {
    var authAPIClient: AuthAPIClient {
        get { self[AuthAPIClientKey.self] }
        set { self[AuthAPIClientKey.self] = newValue }
    }

    private enum AuthAPIClientKey: DependencyKey {
        static let liveValue = AuthAPIClient.live
    }
}

extension AuthAPIClient {
    static let live = AuthAPIClient(
        signIn: { username, password in
            do {
                // 既存のサインイン状態をチェックしてサインアウト
                let currentSession = try await Amplify.Auth.fetchAuthSession()
                if currentSession.isSignedIn {
                    print("🔄 User already signed in, signing out first...")
                    _ = try await Amplify.Auth.signOut()
                }
                
                let signInResult = try await Amplify.Auth.signIn(username: username, password: password)

                guard signInResult.isSignedIn else {
                    throw NSError(domain: "Auth", code: 401, userInfo: [NSLocalizedDescriptionKey: "Sign-in failed"])
                }

                let session = try await Amplify.Auth.fetchAuthSession()

                guard let cognitoSession = session as? AuthCognitoTokensProvider else {
                    throw NSError(domain: "Auth", code: 500, userInfo: [NSLocalizedDescriptionKey: "Not a Cognito session"])
                }

                let tokensResult = cognitoSession.getCognitoTokens()
                let tokens = try tokensResult.get()

                return AuthTokens(
                    idToken: tokens.idToken,
                    accessToken: tokens.accessToken,
                    refreshToken: tokens.refreshToken
                )
            } catch let error as AuthError {
                print("🔴 AuthError details: \(error)")
                print("🔴 AuthError localizedDescription: \(error.localizedDescription)")
                print("🔴 AuthError underlyingError: \(error.underlyingError?.localizedDescription ?? "None")")
                throw error
            } catch {
                print("🔴 Other error: \(error)")
                throw error
            }
        },
        signOut: {
            _ = try await Amplify.Auth.signOut()
        }
    )
}
