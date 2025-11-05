//
//  SelfPomodoroApp.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2024/11/13.
//

import ComposableArchitecture
import SwiftUI
import SwiftData

@main
struct SelfPomodoroApp: App {
    @State private var modelContainer: ModelContainer = SwiftDataStack.makeContainer()
    @State private var dependencies = DependencyValues._current
    @State private var hasInitialized = false

    var body: some Scene {
        WindowGroup {
            MainView(dependencies: dependencies)
                .task {
                    await initializeAppIfNeeded()
                }
        }
        .modelContainer(modelContainer)
    }

    @MainActor
    private func initializeAppIfNeeded() async {
        guard !hasInitialized else { return }

        dependencies.configureAppDependencies(modelContainer: modelContainer)

        let initializer = AppInitializer(
            userRepository: dependencies.userRepository,
            userIdentifierProvider: dependencies.userIdentifier
        )

        await initializer.initialize()
        hasInitialized = true
    }
}
