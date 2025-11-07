//
//  UserConfigRepositoryDependency.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/07.
//

import ComposableArchitecture

private enum UserConfigRepositoryKey: DependencyKey {
    static var liveValue: any UserConfigRepository {
        fatalError("UserConfigRepository has not been configured")
    }
}

extension DependencyValues {
    var userConfigRepository: any UserConfigRepository {
        get { self[UserConfigRepositoryKey.self] }
        set { self[UserConfigRepositoryKey.self] = newValue }
    }
}
