//
//  UserIdentifierDependency.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/05.
//

import ComposableArchitecture

private enum UserIdentifierProviderKey: DependencyKey {
    static let liveValue: () -> String = {
        UserIdentifierProvider.resolve()
    }
}

extension DependencyValues {
    var userIdentifier: () -> String {
        get { self[UserIdentifierProviderKey.self] }
        set { self[UserIdentifierProviderKey.self] = newValue }
    }
}
