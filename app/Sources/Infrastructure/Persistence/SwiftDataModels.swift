//
//  SwiftDataModels.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/05.
//

import Foundation
import SwiftData

@Model
final class UserModel {
    @Attribute(.unique) var id: UUID
    @Attribute(.unique) var identifier: String
    var createdAt: Date
    var updatedAt: Date

    @Relationship(deleteRule: .cascade) var tasks: [TaskModel]
    init(
        id: UUID = UUID(),
        identifier: String,
        createdAt: Date = Date(),
        updatedAt: Date = Date()
    ) {
        self.id = id
        self.identifier = identifier
        self.createdAt = createdAt
        self.updatedAt = updatedAt
        self.tasks = []
    }
}

// userIdentifierのindexは手動で設定する必要がある？
@Model
final class TaskModel {
    @Attribute(.unique) var id: UUID
    var userIdentifier: String
    var detail: String
    var isCompleted: Bool
    var createdAt: Date
    var updatedAt: Date

    @Relationship var user: UserModel?

    init(
        id: UUID = UUID(),
        userIdentifier: String,
        detail: String,
        isCompleted: Bool = false,
        createdAt: Date = Date(),
        updatedAt: Date = Date(),
        user: UserModel? = nil
    ) {
        self.id = id
        self.userIdentifier = userIdentifier
        self.detail = detail
        self.isCompleted = isCompleted
        self.createdAt = createdAt
        self.updatedAt = updatedAt
        self.user = user
    }
}

@Model
final class SessionRecordModel {
    @Attribute(.unique) var id: UUID
    var userIdentifier: String
    var sessionIdentifier: UUID
    var roundCount: Int
    var totalWorkMinutes: Double
    var breakMinutes: Double
    var avgFocusScore: Double
    var isAborted: Bool
    var createdAt: Date
    var updatedAt: Date

    @Relationship var user: UserModel?

    init(
        id: UUID = UUID(),
        userIdentifier: String,
        sessionIdentifier: UUID,
        roundCount: Int,
        totalWorkMinutes: Double,
        breakMinutes: Double,
        avgFocusScore: Double,
        isAborted: Bool = false,
        createdAt: Date = Date(),
        updatedAt: Date = Date(),
        user: UserModel? = nil
    ) {
        self.id = id
        self.userIdentifier = userIdentifier
        self.sessionIdentifier = sessionIdentifier
        self.roundCount = roundCount
        self.totalWorkMinutes = totalWorkMinutes
        self.breakMinutes = breakMinutes
        self.avgFocusScore = avgFocusScore
        self.isAborted = isAborted
        self.createdAt = createdAt
        self.updatedAt = updatedAt
        self.user = user
    }
}

@Model
final class RoundRecordModel {
    @Attribute(.unique) var id: UUID
    var userIdentifier: String
    var sessionIdentifier: UUID
    var workMinutes: Double
    var breakMinutes: Double
    var focusScore: Int?
    var isAborted: Bool
    var createdAt: Date
    var updatedAt: Date

    @Relationship var user: UserModel?

    init(
        id: UUID = UUID(),
        userIdentifier: String,
        sessionIdentifier: UUID,
        workMinutes: Double,
        breakMinutes: Double,
        focusScore: Int?,
        isAborted: Bool = false,
        createdAt: Date = Date(),
        updatedAt: Date = Date(),
        user: UserModel? = nil
    ) {
        self.id = id
        self.userIdentifier = userIdentifier
        self.sessionIdentifier = sessionIdentifier
        self.workMinutes = workMinutes
        self.breakMinutes = breakMinutes
        self.focusScore = focusScore
        self.isAborted = isAborted
        self.createdAt = createdAt
        self.updatedAt = updatedAt
        self.user = user
    }
}

@Model
final class UserConfigRoundHistoryModel {
    @Attribute(.unique) var id: UUID
    var userIdentifier: String
    var recommendedWorkMinutes: Double
    var recommendedBreakMinutes: Double
    var createdAt: Date

    @Relationship var user: UserModel?

    init(
        id: UUID = UUID(),
        userIdentifier: String,
        recommendedWorkMinutes: Double,
        recommendedBreakMinutes: Double,
        createdAt: Date = Date(),
        user: UserModel? = nil
    ) {
        self.id = id
        self.userIdentifier = userIdentifier
        self.recommendedWorkMinutes = recommendedWorkMinutes
        self.recommendedBreakMinutes = recommendedBreakMinutes
        self.createdAt = createdAt
        self.user = user
    }
}

@Model
final class UserConfigSessionHistoryModel {
    @Attribute(.unique) var id: UUID
    var userIdentifier: String
    var recommendedSessionRounds: Int
    var recommendedSessionBreakMinutes: Double
    var createdAt: Date

    @Relationship var user: UserModel?

    init(
        id: UUID = UUID(),
        userIdentifier: String,
        recommendedSessionRounds: Int,
        recommendedSessionBreakMinutes: Double,
        createdAt: Date = Date(),
        user: UserModel? = nil
    ) {
        self.id = id
        self.userIdentifier = userIdentifier
        self.recommendedSessionRounds = recommendedSessionRounds
        self.recommendedSessionBreakMinutes = recommendedSessionBreakMinutes
        self.createdAt = createdAt
        self.user = user
    }
}

@Model
final class UserConfigLatestModel {
    @Attribute(.unique) var id: UUID
    @Attribute(.unique) var userIdentifier: String
    var roundWorkMinutes: Double
    var roundBreakMinutes: Double
    var sessionRounds: Int
    var sessionBreakMinutes: Double
    var updatedAt: Date

    @Relationship var user: UserModel?

    init(
        id: UUID = UUID(),
        userIdentifier: String,
        roundWorkMinutes: Double,
        roundBreakMinutes: Double,
        sessionRounds: Int,
        sessionBreakMinutes: Double,
        updatedAt: Date = Date(),
        user: UserModel? = nil
    ) {
        self.id = id
        self.userIdentifier = userIdentifier
        self.roundWorkMinutes = roundWorkMinutes
        self.roundBreakMinutes = roundBreakMinutes
        self.sessionRounds = sessionRounds
        self.sessionBreakMinutes = sessionBreakMinutes
        self.updatedAt = updatedAt
        self.user = user
    }
}
