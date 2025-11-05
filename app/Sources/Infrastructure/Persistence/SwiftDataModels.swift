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

@Model
final class TaskModel {
    @Attribute(.unique) var id: UUID
    @Attribute(.indexed) var userIdentifier: String
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
