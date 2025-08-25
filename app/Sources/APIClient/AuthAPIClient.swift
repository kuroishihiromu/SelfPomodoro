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
    var signUp: (_ username: String, _ password: String) async throws -> Void
    var confirmSignUp: (_ username: String, _ confirmationCode: String) async throws -> AuthTokens
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
            } catch {
                throw error
            }
        },
        signUp: { username, password in
            do {
                // 既存のサインイン状態をチェックしてサインアウト
                let currentSession = try await Amplify.Auth.fetchAuthSession()
                if currentSession.isSignedIn {
                    _ = try await Amplify.Auth.signOut()
                }
                
                let signUpResult = try await Amplify.Auth.signUp(
                    username: username,
                    password: password
                )
                
                if !signUpResult.isSignUpComplete {
                    // 確認が必要な場合は正常終了（UIで確認コード入力画面を表示）
                    return
                }
            } catch {
                throw error
            }
        },
        confirmSignUp: { username, confirmationCode in
            do {
                let confirmResult = try await Amplify.Auth.confirmSignUp(
                    for: username,
                    confirmationCode: confirmationCode
                )
                
                guard confirmResult.isSignUpComplete else {
                    throw NSError(domain: "Auth", code: 400, userInfo: [NSLocalizedDescriptionKey: "確認が完了しませんでした"])
                }
                
                // 確認完了のみ返す（サインインは別途実行）
                // ダミーのAuthTokensを返す（実際のサインインは確認後に別途実行）
                return AuthTokens(
                    idToken: "confirmed",
                    accessToken: "confirmed", 
                    refreshToken: "confirmed"
                )
            } catch {
                throw error
            }
        },
        signOut: {
            _ = try await Amplify.Auth.signOut()
        }
    )
}
