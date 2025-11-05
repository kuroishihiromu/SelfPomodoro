//
//  AppDependencies.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/05.
//

import SwiftData
import ComposableArchitecture

extension DependencyValues {
    mutating func configureAppDependencies(modelContainer: ModelContainer) {
        let context = ModelContext(modelContainer)
        self.userRepository = SwiftDataUserRepository(context: context)
        self.taskRepository = SwiftDataTaskRepository(context: context)
    }
}
