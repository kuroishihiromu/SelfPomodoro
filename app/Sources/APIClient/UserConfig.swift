//
//  ConfigAPIClient.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/07/03.
//

import Foundation
import Dependencies
import Dependencies

enum userconfigAPIError: Error, Equatable {
    case networkError
    case decodingError
    case unknown
}

struct UserConfigAPIClient {
    var getUserConfig: () async throws -> UserConfigResult
}

extension UserConfigAPIClient {
    static let live = UserConfigAPIClient(
        getUserConfig: {
            // TODO: SwiftDataへの移行が完了したら永続化された最適化設定を参照する
            return UserConfigResult(
                id: UUID(),
                roundWorkTime: 25,
                roundBreakTime: 5,
                sessionRounds: 5,
                sessionBreakTime: 15
            )
        }
    )
}

extension DependencyValues {
    var userConfigAPIClient: UserConfigAPIClient {
        get { self[UserConfigAPIClientKey.self] }
        set { self[UserConfigAPIClientKey.self] = newValue }
    }

    private enum UserConfigAPIClientKey: DependencyKey {
        static let liveValue = UserConfigAPIClient.live
    }
}
