//
//  SessionRecordRepositoryDependency.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/07.
//

import ComposableArchitecture

private enum SessionRecordRepositoryKey: DependencyKey {
    static var liveValue: any SessionRecordRepository {
        fatalError("SessionRecordRepository has not been configured")
    }
}

extension DependencyValues {
    var sessionRecordRepository: any SessionRecordRepository {
        get { self[SessionRecordRepositoryKey.self] }
        set { self[SessionRecordRepositoryKey.self] = newValue }
    }
}
