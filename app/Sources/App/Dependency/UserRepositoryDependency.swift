//
//  UserRepositoryDependency.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/05.
//

import ComposableArchitecture

private enum UserRepositoryKey: DependencyKey {
    static var liveValue: any UserRepository {
        fatalError("UserRepository live value has not been configured")
    }
}

extension DependencyValues {
    var userRepository: any UserRepository {
        get { self[UserRepositoryKey.self] }
        set { self[UserRepositoryKey.self] = newValue }
    }
}
