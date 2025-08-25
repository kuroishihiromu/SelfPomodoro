//
//  AuthFeature.swift
//  SelfPomodoro
//

//  Created by 黒石陽夢 on 2025/06/16.
//

import ComposableArchitecture
import Foundation

@Reducer
struct AuthFeature: Reducer {
    @ObservableState
    struct State: Equatable {
        var email = ""
        var password = ""
        var confirmPassword = ""
        var otpCode = ""
        var isAgreed = false
        var errorMessage: String?
        var tokens: AuthTokens?
        var isLoggedIn = false
        var showOTPInput = false
        var otpSent = false
    }

    enum Action: BindableAction {
        case binding(BindingAction<State>)
        case tappedLogin
        case tappedSignUp
        case tappedSignOut
        case tappedConfirmOTP
        case loginResponse(Result<AuthTokens, Error>)
        case signUpResponse(Result<Void, Error>)
        case confirmOTPResponse(Result<AuthTokens, Error>)
        case signOutResponse(Result<Void, Error>)
    }

    @Dependency(\.authAPIClient) var authAPIClient

    var body: some ReducerOf<Self> {
        BindingReducer()

        Reduce { state, action in
            switch action {
            case .tappedLogin:
                guard state.isAgreed else {
                    state.errorMessage = L10n.Auth.termsRequired
                    return .none
                }

                return .run { [email = state.email, password = state.password] send in
                    do {
                        let tokens = try await authAPIClient.signIn(email, password)
                        await send(.loginResponse(.success(tokens)))
                    } catch {
                        await send(.loginResponse(.failure(error)))
                    }
                }

            case .tappedSignUp:
                guard state.isAgreed else {
                    state.errorMessage = L10n.Auth.termsRequired
                    return .none
                }
                
                guard state.password == state.confirmPassword else {
                    state.errorMessage = L10n.Auth.passwordMismatch
                    return .none
                }

                return .run { [email = state.email, password = state.password] send in
                    do {
                        try await authAPIClient.signUp(email, password)
                        await send(.signUpResponse(.success(())))
                    } catch {
                        await send(.signUpResponse(.failure(error)))
                    }
                }

            case let .loginResponse(.success(tokens)):
                state.tokens = tokens
                state.isLoggedIn = true
                state.errorMessage = nil
                
                return .none

            case let .loginResponse(.failure(error)):
                state.errorMessage = error.localizedDescription
                return .none

            case .signUpResponse(.success):
                state.showOTPInput = true
                state.otpSent = true
                state.errorMessage = nil
                
                return .none

            case let .signUpResponse(.failure(error)):
                state.errorMessage = error.localizedDescription
                return .none

            case .tappedConfirmOTP:
                guard !state.otpCode.isEmpty else {
                    state.errorMessage = "確認コードを入力してください"
                    return .none
                }

                return .run { [email = state.email, otpCode = state.otpCode] send in
                    do {
                        let tokens = try await authAPIClient.confirmSignUp(email, otpCode)
                        await send(.confirmOTPResponse(.success(tokens)))
                    } catch {
                        await send(.confirmOTPResponse(.failure(error)))
                    }
                }

            case let .confirmOTPResponse(.success(tokens)):
                if tokens.idToken == "confirmed" {
                    // 確認完了、サインイン画面に戻る
                    state.showOTPInput = false
                    state.otpSent = false
                    state.otpCode = ""
                    state.errorMessage = "確認が完了しました。サインインしてください。"
                } else {
                    // 直接ログイン
                    state.tokens = tokens
                    state.isLoggedIn = true
                    state.errorMessage = nil
                }
                return .none

            case let .confirmOTPResponse(.failure(error)):
                state.errorMessage = error.localizedDescription
                return .none

            case .tappedSignOut:
                return .run { send in
                    do {
                        try await authAPIClient.signOut()
                        await send(.signOutResponse(.success(())))
                    } catch {
                        await send(.signOutResponse(.failure(error)))
                    }
                }

            case .signOutResponse(.success):
                state.tokens = nil
                state.isLoggedIn = false
                state.errorMessage = nil
                state.email = ""
                state.password = ""
                state.confirmPassword = ""
                state.isAgreed = false
                return .none

            case let .signOutResponse(.failure(error)):
                state.errorMessage = error.localizedDescription
                return .none

            case .binding:
                return .none
            }
        }
    }
}
