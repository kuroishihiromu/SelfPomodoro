//
//  SelfPomodoroApp.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2024/11/13.
//

import ComposableArchitecture
import SwiftUI

// @main
// struct SelfPomodoroApp: App {
//    @State private var isAuthenticated = false
//    var body: some Scene {
//        WindowGroup {
//            AuthScreenView()
//                .environment(\.isAuthenticated, isAuthenticated)
//                .task {
//                    for await state in supabase.auth.authStateChanges {
//                        if [.initialSession, .signedIn, .signedOut].contains(state.event) {
//                            isAuthenticated = state.session != nil
//                        }
//                    }
//                }
//        }
//    }
// }

@main
struct SelfPomodoroApp: App {
    var body: some Scene {
        WindowGroup {
            MainView(
                store: Store(initialState: TabButtonFeature.State()) {
                    TabButtonFeature()
                }
            )
        }
    }
}
