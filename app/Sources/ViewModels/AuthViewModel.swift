//
//  AuthViewModel.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/02/27.
//

import Foundation
import Supabase
import SwiftUI

@MainActor
class AuthViewModel: ObservableObject {
    @Published var email = ""
    @Published var password = ""
    @Published var isLoading = false
    @Published var result: Result<Void, Error>?
    @Published var isAuthenticated = false // 認証状態を管理

    init() {
        Task {
            await checkSession()
            await observeAuthChanges()
        }
    }

    func checkSession() async {
        if let session = try? await supabase.auth.session {
            DispatchQueue.main.async {
                self.isAuthenticated = session != nil
            }
        }
    }

    func signIn() async {
        isLoading = true
        defer { isLoading = false }

        do {
            try await supabase.auth.signIn(email: email, password: password)
            DispatchQueue.main.async {
                self.isAuthenticated = true
            }
            result = .success(())
        } catch {
            result = .failure(error)
        }
    }

    func signUp() async {
        isLoading = true
        defer { isLoading = false }

        do {
            try await supabase.auth.signUp(email: email, password: password)
            DispatchQueue.main.async {
                self.isAuthenticated = true
            }
            result = .success(())
        } catch {
            print("SignUp Error:", error.localizedDescription)
            result = .failure(error)
        }
    }

    func signOut() async {
        isLoading = true
        defer { isLoading = false }

        do {
            try await supabase.auth.signOut()
            DispatchQueue.main.async {
                self.isAuthenticated = false // ✅ ログアウト時に状態を更新
            }
            result = .success(())
        } catch {
            result = .failure(error)
        }
    }

    func observeAuthChanges() async {
        for await state in supabase.auth.authStateChanges {
            print("Auth State Changed: \(state.event), Session: \(String(describing: state.session))")

            DispatchQueue.main.async {
                if state.event == .initialSession {
                    print("Skipping initial session check")
                    return
                }
                self.isAuthenticated = state.session != nil
                print("Updated isAuthenticated: \(self.isAuthenticated)")
            }
        }
    }
}
