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
    @State private var modelContainer: ModelContainer
    @State private var dependencies: DependencyValues
    @State private var hasInitialized = false

    init() {
        let container = SwiftDataStack.makeContainer()
        var deps = DependencyValues._current
        deps.configureAppDependencies(modelContainer: container)
        _modelContainer = State(initialValue: container)
        _dependencies = State(initialValue: deps)
    }

    var body: some Scene {
        WindowGroup {
            Group {
                if hasInitialized {
                    MainView(dependencies: dependencies)
                } else {
                    ProgressView("Loading...")
                }
            }
            .task {
                await initializeAppIfNeeded()
            }
        }
        .modelContainer(modelContainer)
    }

    @MainActor
    private func initializeAppIfNeeded() async {
        guard !hasInitialized else { return }

        let initializer = AppInitializer(
            userRepository: dependencies.userRepository,
            userIdentifierProvider: dependencies.userIdentifier
        )

        await initializer.initialize()
        hasInitialized = true
    }
}
