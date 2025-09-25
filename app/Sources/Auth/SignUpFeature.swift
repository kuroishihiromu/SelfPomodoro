//
//  SignUpFeature.swift
//  SelfPomodoro
//
//  Created by kuroishi hiromu on 2025/08/25.
//

import ComposableArchitecture
import Foundation

@Reducer
struct SignUpFeature: Reducer {
    @ObservableState
    struct State: Equatable {
        var email = ""
        var displayname = ""
        var password = ""
        var confirmPassword = ""
        var otpCode = ""
        var isAgreed = false
        var errorMessage: String?
        var showOTPInput = false
        var otpSent = false
    }

    enum Action: BindableAction {
        case binding(BindingAction<State>)
        case tappedSignUp
        case tappedConfirmOTP
        case signUpResponse(Result<Void, Error>)
        case confirmOTPResponse(Result<AuthTokens, Error>)
        case delegate(Delegate)
        
        enum Delegate {
            case signUpCompleted
            case userSignedIn(AuthTokens)
        }
    }

    @Dependency(\.authAPIClient) var authAPIClient

    var body: some ReducerOf<Self> {
        BindingReducer()

        Reduce { state, action in
            switch action {
            case .tappedSignUp:
                print("🔵 DEBUG: tappedSignUp action called")
                print("🔵 DEBUG: email=\(state.email), password length=\(state.password.count)")
                print("🔵 DEBUG: isAgreed=\(state.isAgreed)")
                
                // バリデーション
                guard !state.displayname.isEmpty else {
                    state.errorMessage = L10n.Validation.displayNameRequired
                    return .none
                }
                
                guard state.displayname.count >= 2 else {
                    state.errorMessage = L10n.Validation.displayNameTooShort
                    return .none
                }
                
                guard state.displayname.count <= 50 else {
                    state.errorMessage = L10n.Validation.displayNameTooLong
                    return .none
                }
                
                guard !state.email.isEmpty else {
                    state.errorMessage = L10n.Validation.emailRequired
                    return .none
                }
                
                guard state.email.contains("@") && state.email.contains(".") else {
                    state.errorMessage = L10n.Validation.emailInvalid
                    return .none
                }
                
                guard !state.password.isEmpty else {
                    state.errorMessage = L10n.Validation.passwordRequired
                    return .none
                }
                
                guard state.password.count >= 8 else {
                    state.errorMessage = L10n.Validation.passwordTooShort
                    return .none
                }
                
                // パスワード複雑性チェック（英大文字・英小文字・数字を少なくとも1つ含む）
                let hasLowercase = state.password.contains { $0.isLowercase }
                let hasUppercase = state.password.contains { $0.isUppercase }
                let hasDigit = state.password.contains { $0.isNumber }
                
                guard hasLowercase && hasUppercase && hasDigit else {
                    state.errorMessage = L10n.Validation.passwordWeak
                    return .none
                }
                
                guard !state.confirmPassword.isEmpty else {
                    state.errorMessage = L10n.Validation.confirmPasswordRequired
                    return .none
                }
                
                guard state.password == state.confirmPassword else {
                    state.errorMessage = L10n.Validation.passwordMismatch
                    return .none
                }
                
                guard state.isAgreed else {
                    print("🔴 DEBUG: Terms not agreed")
                    state.errorMessage = L10n.Auth.termsRequired
                    return .none
                }

                print("🔵 DEBUG: Starting signUp API call")
                return .run { [email = state.email, displayname = state.displayname, password = state.password] send in
                    do {
                        print("🔵 DEBUG: Calling authAPIClient.signUp with email: \(email), displayname: \(displayname)")
                        try await authAPIClient.signUp(email, password, displayname)
                        print("🟢 DEBUG: signUp API call successful")
                        await send(.signUpResponse(.success(())))
                    } catch {
                        print("🔴 DEBUG: signUp API call failed: \(error)")
                        await send(.signUpResponse(.failure(error)))
                    }
                }

            case .signUpResponse(.success):
                print("🟢 DEBUG: signUpResponse success - showing OTP input")
                state.showOTPInput = true
                state.otpSent = true
                state.errorMessage = nil
                
                return .none

            case let .signUpResponse(.failure(error)):
                print("🔴 DEBUG: signUpResponse failure: \(error.localizedDescription)")
                state.errorMessage = ErrorMessageHelper.localizedAuthError(error)
                return .none

            case .tappedConfirmOTP:
                guard !state.otpCode.isEmpty else {
                    state.errorMessage = L10n.Validation.otpRequired
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
                    state.errorMessage = L10n.Success.confirmationCompleted
                    return .send(.delegate(.signUpCompleted))
                } else {
                    // 直接ログイン
                    return .send(.delegate(.userSignedIn(tokens)))
                }

            case let .confirmOTPResponse(.failure(error)):
                state.errorMessage = ErrorMessageHelper.localizedAuthError(error)
                return .none

            case .binding:
                return .none
                
            case .delegate:
                return .none
            }
        }
    }
}
