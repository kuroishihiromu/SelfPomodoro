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
        let userRepository = SwiftDataUserRepository(context: context)
        let taskRepository = SwiftDataTaskRepository(context: context)
        let sessionRecordRepository = SwiftDataSessionRecordRepository(context: context)
        let roundRecordRepository = SwiftDataRoundRecordRepository(context: context)
        let userConfigRepository = SwiftDataUserConfigRepository(context: context)

        self.userRepository = userRepository
        self.taskRepository = taskRepository
        self.sessionRecordRepository = sessionRecordRepository
        self.roundRecordRepository = roundRecordRepository
        self.userConfigRepository = userConfigRepository
    }

    @MainActor
    static func makeConfiguredDependencies(modelContainer: ModelContainer) -> DependencyValues {
        var dependencies = DependencyValues._current
        dependencies.configureAppDependencies(modelContainer: modelContainer)
        return dependencies
    }
}
