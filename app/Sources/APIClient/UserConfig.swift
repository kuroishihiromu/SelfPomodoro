//
//  ConfigAPIClient.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/07/03.
//

import Foundation
import Dependencies
import Amplify

struct UserConfigAPIClient {
    var getUserConfig: () async throws -> UserConfigResult
}

extension UserConfigAPIClient {
    static let live = UserConfigAPIClient(
        getUserConfig: {
            let request = RESTRequest(
                apiName: "selfpomodoro",
                path: "/dev/api/v1/user-config",
                headers: [
                    "Content-Type": "application/json"
                ]
            )

            let data = try await Amplify.API.get(request: request)
            print("Get Config → \(String(data: data, encoding: .utf8) ?? "Invalid UTF-8")")
            return try APIFormatters.jsonDecoder.decode(UserConfigResult.self, from: data)
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
