//
//  SelfPomodoroApp.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2024/11/13.
//

import SwiftUI
import Amplify
import AWSCognitoAuthPlugin
import ComposableArchitecture
import AWSPluginsCore
import AWSAPIPlugin

@main
struct SelfPomodoroApp: App {
    @State private var isConfigured = false
    @State private var isSignedIn: Bool? = nil
    @State private var tokens: AuthTokens? = nil

    var body: some Scene {
        WindowGroup {
            Group {
                if !isConfigured {
                    ProgressView("Initializing...")
                } else if isSignedIn == true, let tokens {
                    MainView(token: tokens)
                } else {
                    AuthScreenView(
                        store: Store(
                            initialState: AuthFeature.State(),
                            reducer: { AuthFeature() }
                        )
                    )
                }
            }
            .task {
                await initializeApp()
            }
        }
    }

    @MainActor
    private func initializeApp() async {
        do {
            try Amplify.add(plugin: AWSCognitoAuthPlugin())
            try Amplify.add(plugin: AWSAPIPlugin())
            try Amplify.configure()
            print("✅ Amplify configured")
            isConfigured = true

            let session = try await Amplify.Auth.fetchAuthSession()
            if session.isSignedIn,
               let provider = session as? AuthCognitoTokensProvider {
                let result = provider.getCognitoTokens()
                let token = try result.get()
                tokens = AuthTokens(
                    idToken: token.idToken,
                    accessToken: token.accessToken,
                    refreshToken: token.refreshToken
                )
                isSignedIn = true
            } else {
                isSignedIn = false
            }
        } catch {
            print("❌ Initialization failed: \(error)")
            isSignedIn = false
        }
    }
}
