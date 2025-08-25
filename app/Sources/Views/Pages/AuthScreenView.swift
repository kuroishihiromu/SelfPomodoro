//
//  AuthScreenView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/02/27.
//

import SwiftUI
import ComposableArchitecture

struct AuthScreenView: View {
    @Bindable var store: StoreOf<AuthFeature>

    var body: some View {
        if store.isLoggedIn {
            MainView(token: store.tokens!, authStore: store)
        } else {
            NavigationStack(path: $store.scope(state: \.path, action: \.path)) {
                VStack(spacing: 20) {
                    Text(L10n.Auth.createAccount)
                        .font(.system(size: 30, weight: .bold))
                    Text(L10n.Auth.description)
                        .multilineTextAlignment(.center)
                        .frame(maxWidth: .infinity, alignment: .center)
                        .padding(.bottom, 50)
                    NormalTextField(
                        placeholder: L10n.Auth.emailPlaceholder,
                        icon: Image(.mail),
                        width: 350,
                        height: 44,
                        text: $store.email
                    )
                    PasswordTextField(
                        placeholder: L10n.Auth.passwordPlaceholder,
                        icon: Image(.key),
                        width: 350,
                        height: 44,
                        text: $store.password
                    )
                    
                    // エラーメッセージ（統一表示）
                    if let message = store.errorMessage {
                        Text(message)
                            .foregroundColor(.red)
                            .padding()
                            .multilineTextAlignment(.center)
                    }
                    
                    NormalButton(text: L10n.Auth.loginButton, bgColor: ColorTheme.navy, fontColor: ColorTheme.white, width: 350, height: 50, action: {store.send(.tappedLogin)})
                                        
                    Button(action: { store.send(.navigateToSignUp) }) {
                        Text(L10n.Auth.signupLink)
                            .foregroundColor(ColorTheme.black)
                            .underline()
                    }
                    .padding(.top, 10)

                }
                .frame(width: 350)
            } destination: { pathStore in
                switch pathStore.state {
                case .signUp:
                    SignUpScreenView(store: store)
                }
            }
        }
    }
}

#Preview {
    AuthScreenView(
        store: Store(
            initialState: AuthFeature.State(),
            reducer: { AuthFeature() }
        )
    )
}
