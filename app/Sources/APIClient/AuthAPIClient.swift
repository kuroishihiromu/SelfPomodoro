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
    var signIn: (_ email: String, _ password: String) async throws -> AuthTokens
    var signUp: (_ email: String, _ password: String, _ name: String) async throws -> Void
    var confirmSignUp: (_ email: String, _ confirmationCode: String) async throws -> AuthTokens
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
        signIn: { email, password in
            do {
                // 既存のサインイン状態をチェックしてサインアウト
                let currentSession = try await Amplify.Auth.fetchAuthSession()
                if currentSession.isSignedIn {
                    _ = try await Amplify.Auth.signOut()
                }
                
                let signInResult = try await Amplify.Auth.signIn(username: email, password: password)

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
        signUp: { email, password, name in
            do {
                print("🔵 DEBUG: AuthAPIClient.signUp called with email: \(email), name: \(name)")
                
                // 既存のサインイン状態をチェックしてサインアウト
                let currentSession = try await Amplify.Auth.fetchAuthSession()
                print("🔵 DEBUG: Current session isSignedIn: \(currentSession.isSignedIn)")
                
                if currentSession.isSignedIn {
                    print("🔵 DEBUG: Signing out existing session")
                    _ = try await Amplify.Auth.signOut()
                }
                
                print("🔵 DEBUG: Calling Amplify.Auth.signUp with user attributes")
                let userAttributes = [
                    AuthUserAttribute(.name, value: name),
                    AuthUserAttribute(.email, value: email)
                ]
                
                let signUpResult = try await Amplify.Auth.signUp(
                    username: email,
                    password: password,
                    options: AuthSignUpRequest.Options(
                        userAttributes: userAttributes
                    )
                )
                
                print("🔵 DEBUG: SignUp result - isSignUpComplete: \(signUpResult.isSignUpComplete)")
                print("🔵 DEBUG: SignUp result - nextStep: \(signUpResult.nextStep)")
                
                if !signUpResult.isSignUpComplete {
                    print("🟢 DEBUG: Confirmation required - email should be sent")
                    // 確認が必要な場合は正常終了（UIで確認コード入力画面を表示）
                    return
                }
                
                print("🟢 DEBUG: SignUp completed without confirmation")
            } catch {
                print("🔴 DEBUG: AuthAPIClient.signUp error: \(error)")
                throw error
            }
        },
        confirmSignUp: { email, confirmationCode in
            do {
                let confirmResult = try await Amplify.Auth.confirmSignUp(
                    for: email,
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
