//
//  AppInitializer.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/05.
//

import Foundation

@MainActor
final class AppInitializer {
    private let userRepository: UserRepository
    private let userIdentifierProvider: () -> String

    init(
        userRepository: UserRepository,
        userIdentifierProvider: @escaping () -> String
    ) {
        self.userRepository = userRepository
        self.userIdentifierProvider = userIdentifierProvider
    }

    func initialize() async {
        let identifier = userIdentifierProvider()
        do {
            if try await userRepository.fetchUser(by: identifier) == nil {
                _ = try await userRepository.createUser(identifier: identifier)
            }
        } catch {
            // TODO: エラーロギング戦略を決めたら差し替え
            print("⚠️ AppInitializer failed to ensure user: \(error)")
        }
    }
}
