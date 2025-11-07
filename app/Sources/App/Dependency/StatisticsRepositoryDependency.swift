//
//  StatisticsRepositoryDependency.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/07.
//

import ComposableArchitecture

private enum StatisticsRepositoryKey: DependencyKey {
    static var liveValue: any StatisticsRepository {
        fatalError("StatisticsRepository has not been configured")
    }
}

extension DependencyValues {
    var statisticsRepository: any StatisticsRepository {
        get { self[StatisticsRepositoryKey.self] }
        set { self[StatisticsRepositoryKey.self] = newValue }
    }
}
