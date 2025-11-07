//
//  SwiftDataSessionRecordRepository.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/07.
//

import Foundation
import SwiftData

@MainActor
final class SwiftDataSessionRecordRepository: SessionRecordRepository {
    private let context: ModelContext

    init(context: ModelContext) {
        self.context = context
    }

    func save(_ record: SessionRecord) async throws {
        let user = try fetchUserModel(by: record.userIdentifier)
        let model = SessionRecordModel(
            id: record.id,
            userIdentifier: record.userIdentifier,
            sessionIdentifier: record.sessionIdentifier,
            roundCount: record.roundCount,
            totalWorkMinutes: record.totalWorkMinutes,
            breakMinutes: record.breakMinutes,
            avgFocusScore: record.averageFocusScore,
            isAborted: record.isAborted,
            createdAt: record.createdAt,
            updatedAt: record.updatedAt,
            user: user
        )
        context.insert(model)
        try context.save()
        print("🧾 SessionRecord saved id=\(model.id) aborted=\(model.isAborted)")
    }

    func fetchRecent(for userIdentifier: String, limit: Int?) async throws -> [SessionRecord] {
        var descriptor = FetchDescriptor<SessionRecordModel>(
            predicate: #Predicate { $0.userIdentifier == userIdentifier },
            sortBy: [SortDescriptor(\.createdAt, order: .reverse)]
        )
        if let limit {
            descriptor.fetchLimit = limit
        }
        let models = try context.fetch(descriptor)
        return try models.map { try map($0) }
    }

    // MARK: - Helpers

    private func fetchUserModel(by identifier: String) throws -> UserModel {
        var descriptor = FetchDescriptor<UserModel>(
            predicate: #Predicate { $0.identifier == identifier }
        )
        descriptor.fetchLimit = 1
        guard let model = try context.fetch(descriptor).first else {
            throw RepositoryError.userNotFound
        }
        return model
    }

    private func map(_ model: SessionRecordModel) throws -> SessionRecord {
        return SessionRecord(
            id: model.id,
            userIdentifier: model.userIdentifier,
            sessionIdentifier: model.sessionIdentifier,
            roundCount: model.roundCount,
            totalWorkMinutes: model.totalWorkMinutes,
            breakMinutes: model.breakMinutes,
            averageFocusScore: model.avgFocusScore,
            isAborted: model.isAborted,
            createdAt: model.createdAt,
            updatedAt: model.updatedAt
        )
    }
}
