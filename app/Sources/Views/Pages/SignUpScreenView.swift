//
//  SignUpScreenView.swift
//  SelfPomodoro
//
//  Created by kuroishi hiromu on 2025/08/25.
//

import SwiftUI
import ComposableArchitecture

struct SignUpScreenView: View {
    @Bindable var store: StoreOf<AuthFeature>

    var body: some View {
        if store.showOTPInput {
            // OTP入力画面
            VStack(spacing: 20) {
                Text(L10n.Auth.confirmTitle)
                    .font(.system(size: 30, weight: .bold))
                Text(L10n.Auth.confirmDescription)
                    .multilineTextAlignment(.center)
                    .frame(maxWidth: .infinity, alignment: .center)
                    .padding(.bottom, 50)
                
                NormalTextField(
                    placeholder: L10n.Auth.otpPlaceholder,
                    icon: Image(.mail),
                    width: 350,
                    height: 44,
                    text: $store.otpCode
                )
                
                // エラーメッセージ
                if let message = store.errorMessage {
                    Text(message)
                        .foregroundColor(.red)
                        .padding()
                        .multilineTextAlignment(.center)
                }
                
                NormalButton(
                    text: L10n.Auth.confirmButton,
                    bgColor: ColorTheme.navy,
                    fontColor: ColorTheme.white,
                    width: 350,
                    height: 50,
                    action: { store.send(.tappedConfirmOTP) }
                )
                
                Spacer()
            }
            .frame(width: 350)
            .navigationTitle(L10n.Auth.confirmTitle)
            .navigationBarTitleDisplayMode(.inline)
        } else {
            // サインアップ画面
            VStack(spacing: 20) {
                Text(L10n.Auth.signupDescription)
                    .multilineTextAlignment(.center)
                    .frame(maxWidth: .infinity, alignment: .center)
                    .padding(.bottom, 50)
                
                NormalTextField(
                    placeholder: L10n.Auth.namePlaceholder,
                    icon: Image(systemName: "person"),
                    width: 350,
                    height: 44,
                    text: $store.displayname
                )
                
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
                
                PasswordTextField(
                    placeholder: L10n.Auth.confirmPasswordPlaceholder,
                    icon: Image(.key),
                    width: 350,
                    height: 44,
                    text: $store.confirmPassword
                )
                
                // 利用規約同意チェックボックス
                HStack {
                    CheckboxView(isChecked: $store.isAgreed)
                    Text(L10n.Auth.agreeTerms)
                        .font(.system(size: 14))
                    Spacer()
                }
                .padding(.horizontal, 20)
                
                // エラーメッセージ
                if let message = store.errorMessage {
                    Text(message)
                        .foregroundColor(.red)
                        .padding()
                        .multilineTextAlignment(.center)
                }
                
                NormalButton(
                    text: L10n.Auth.signupButton,
                    bgColor: ColorTheme.navy,
                    fontColor: ColorTheme.white,
                    width: 350,
                    height: 50,
                    action: { store.send(.tappedSignUp) }
                )
                
                Spacer()
            }
            .frame(width: 350)
            .navigationTitle(L10n.Auth.createAccount)
            .navigationBarTitleDisplayMode(.inline)
        }
    }
}

#Preview {
    NavigationStack {
        SignUpScreenView(
            store: Store(
                initialState: AuthFeature.State(),
                reducer: { AuthFeature() }
            )
        )
    }
}
