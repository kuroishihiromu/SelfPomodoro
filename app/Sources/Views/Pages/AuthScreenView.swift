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
            MainView(token: store.tokens!)
        } else {
            VStack(spacing: 20) {
                Text("Create Account")
                    .font(.system(size: 30, weight: .bold))
                Text("Create your account to start managing your time effectively")
                    .multilineTextAlignment(.center)
                    .frame(maxWidth: .infinity, alignment: .center)
                    .padding(.bottom, 50)
                NormalTextField(
                    placeholder: "Your email address",
                    icon: Image(.mail),
                    width: 350,
                    height: 44,
                    text: $store.email
                )
                PasswordTextField(
                    placeholder: "Enter your password",
                    icon: Image(.key),
                    width: 350,
                    height: 44,
                    text: $store.password
                )
                HStack{
                    CheckboxView(isChecked: $store.isAgreed)
                    Text("I agree with Terms & Conditions")
                    Spacer()
                }
                .frame(width: 350)
                NormalButton(text: "Log in", bgColor: ColorTheme.navy, fontColor: ColorTheme.white, width: 350, height: 50, action: {store.send(.tappedLogin)})
                if let message = store.errorMessage {
                    Text(message)
                        .foregroundColor(.red)
                        .padding()
                }
                NormalButton(text: "Sign Up", bgColor: ColorTheme.navy, fontColor: ColorTheme.white, width: 350, height: 50, action: {store.send(.tappedLogin)})
                if let message = store.errorMessage {
                    Text(message)
                        .foregroundColor(.red)
                        .padding()
                }

            }
            .frame(width: 350)
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
