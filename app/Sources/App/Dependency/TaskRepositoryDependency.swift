//
//  TaskRepositoryDependency.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/05.
//

import ComposableArchitecture

private enum TaskRepositoryKey: DependencyKey {
    static var liveValue: any TaskRepository {
        fatalError("TaskRepository live value has not been configured")
    }
}

extension DependencyValues {
    var taskRepository: any TaskRepository {
        get { self[TaskRepositoryKey.self] }
        set { self[TaskRepositoryKey.self] = newValue }
    }
}
