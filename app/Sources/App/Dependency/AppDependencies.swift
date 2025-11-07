//
//  AppDependencies.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/07.
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
        let statisticsRepository = SwiftDataStatisticsRepository(context: context)

        self.userRepository = userRepository
        self.taskRepository = taskRepository
        self.sessionRecordRepository = sessionRecordRepository
        #if DEBUG
        self.roundRecordRepository = DebugRoundRecordRepository(base: roundRecordRepository)
        #else
        self.roundRecordRepository = roundRecordRepository
        #endif
        self.userConfigRepository = userConfigRepository
        self.statisticsRepository = statisticsRepository
    }

    @MainActor
    static func makeConfiguredDependencies(modelContainer: ModelContainer) -> DependencyValues {
        var dependencies = DependencyValues._current
        dependencies.configureAppDependencies(modelContainer: modelContainer)
        return dependencies
    }
}
