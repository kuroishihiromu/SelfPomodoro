//
//  AppDependencies.swift
//  SelfPomodoro
//
//  Created by Codex on 2025/02/15.
//

import SwiftData
import ComposableArchitecture

extension DependencyValues {
    @MainActor
    mutating func configureAppDependencies(modelContainer: ModelContainer) {
        let context = ModelContext(modelContainer)
        self.userRepository = SwiftDataUserRepository(context: context)
        self.taskRepository = SwiftDataTaskRepository(context: context)
    }

    @MainActor
    static func makeConfiguredDependencies(modelContainer: ModelContainer) -> DependencyValues {
        var dependencies = DependencyValues._current
        dependencies.configureAppDependencies(modelContainer: modelContainer)
        return dependencies
    }
}
