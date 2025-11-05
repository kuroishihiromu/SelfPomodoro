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
    @State private var hasInitialized = false

    var body: some Scene {
        WindowGroup {
            MainView()
                .task {
                    await initializeAppIfNeeded()
                }
        }
        .modelContainer(modelContainer)
    }

    @MainActor
    private func initializeAppIfNeeded() async {
        guard !hasInitialized else { return }

        var dependencies = DependencyValues._current
        dependencies.configureAppDependencies(modelContainer: modelContainer)
        DependencyValues._current = dependencies

        let userRepository = dependencies.userRepository
        let initializer = AppInitializer(
            userRepository: userRepository,
            userIdentifierProvider: UserIdentifierProvider.resolve
        )

        await initializer.initialize()
        hasInitialized = true
    }
}
