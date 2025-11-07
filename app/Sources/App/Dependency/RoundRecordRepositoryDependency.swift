//
//  RoundRecordRepositoryDependency.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/07.
//

import ComposableArchitecture

private enum RoundRecordRepositoryKey: DependencyKey {
    static var liveValue: any RoundRecordRepository {
        fatalError("RoundRecordRepository has not been configured")
    }
}

extension DependencyValues {
    var roundRecordRepository: any RoundRecordRepository {
        get { self[RoundRecordRepositoryKey.self] }
        set { self[RoundRecordRepositoryKey.self] = newValue }
    }
}
