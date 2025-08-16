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
        var isAgreed = false
        var errorMessage: String?
        var tokens: AuthTokens?
        var isLoggedIn = false
    }

    enum Action: BindableAction {
        case binding(BindingAction<State>)
        case tappedLogin
        case loginResponse(Result<AuthTokens, Error>)
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

            case let .loginResponse(.success(tokens)):
                state.tokens = tokens
                state.isLoggedIn = true
                state.errorMessage = nil
                
                // トークンをTokenStorageに保存
                TokenStorage.shared.setToken(tokens)
                
                return .none

            case let .loginResponse(.failure(error)):
                state.errorMessage = error.localizedDescription
                return .none

            case .binding:
                return .none
            }
        }
    }
}
