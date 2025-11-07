//
//  UserConfig.swift
//  SelfPomodoro
//
//  Created by Codex on 2025/02/15.
//

import Foundation

struct UserConfig: Identifiable, Equatable {
    let id: UUID
    let userIdentifier: String
    var roundWorkMinutes: Double
    var roundBreakMinutes: Double
    var sessionRounds: Int
    var sessionBreakMinutes: Double
    var updatedAt: Date
}

struct UserConfigRoundHistoryEntry: Identifiable, Equatable {
    let id: UUID
    let userIdentifier: String
    var recommendedWorkMinutes: Double
    var recommendedBreakMinutes: Double
    var createdAt: Date
}

extension UserConfig {
    static func `default`(for userIdentifier: String = "") -> UserConfig {
        UserConfig(
            id: UUID(),
            userIdentifier: userIdentifier,
            roundWorkMinutes: 25,
            roundBreakMinutes: 5,
            sessionRounds: 5,
            sessionBreakMinutes: 15,
            updatedAt: Date()
        )
    }
}

struct UserConfigSessionHistoryEntry: Identifiable, Equatable {
    let id: UUID
    let userIdentifier: String
    var recommendedSessionRounds: Int
    var recommendedSessionBreakMinutes: Double
    var createdAt: Date
}
