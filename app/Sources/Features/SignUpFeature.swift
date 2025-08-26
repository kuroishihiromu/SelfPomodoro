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
        var name = ""
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
                
                guard state.isAgreed else {
                    print("🔴 DEBUG: Terms not agreed")
                    state.errorMessage = L10n.Auth.termsRequired
                    return .none
                }
                
                guard state.password == state.confirmPassword else {
                    print("🔴 DEBUG: Password mismatch")
                    state.errorMessage = L10n.Auth.passwordMismatch
                    return .none
                }

                print("🔵 DEBUG: Starting signUp API call")
                return .run { [email = state.email, name = state.name, password = state.password] send in
                    do {
                        print("🔵 DEBUG: Calling authAPIClient.signUp with email: \(email), name: \(name)")
                        try await authAPIClient.signUp(email, password, name)
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
                    return .send(.delegate(.signUpCompleted))
                } else {
                    // 直接ログイン
                    return .send(.delegate(.userSignedIn(tokens)))
                }

            case let .confirmOTPResponse(.failure(error)):
                state.errorMessage = error.localizedDescription
                return .none

            case .binding:
                return .none
                
            case .delegate:
                return .none
            }
        }
    }
}
