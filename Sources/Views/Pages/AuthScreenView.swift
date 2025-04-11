//
//  AuthView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/02/27.
//

import SwiftUI

struct AuthScreenView: View {
    @StateObject private var viewModel = AuthViewModel()
    @Environment(\.isAuthenticated) private var isAuthenticated
    
    var body: some View {
        Group {
            if viewModel.isAuthenticated {
                HomeView() // 認証成功時に HomeView へ遷移
            } else {
                VStack {
                    TextField("Email", text: $viewModel.email)
                        .textFieldStyle(RoundedBorderTextFieldStyle())
                        .padding()

                    SecureField("Password", text: $viewModel.password)
                        .textFieldStyle(RoundedBorderTextFieldStyle())
                        .padding()

                    Button("Sign In") {
                        Task {
                            await viewModel.signIn()
                        }
                    }
                    .buttonStyle(.borderedProminent)
                    .padding()

                    Button("Sign Up") {
                        Task {
                            await viewModel.signUp()
                        }
                    }
                    .buttonStyle(.bordered)
                    .padding()
                }
            }
        }
        .onAppear {
            Task {
                await viewModel.checkSession()
            }
        }
    }
}
