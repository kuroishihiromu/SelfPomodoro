//
//  SwiftDataRoundRecordRepository.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/07.
//

import Foundation
import SwiftData

@MainActor
final class SwiftDataRoundRecordRepository: RoundRecordRepository {
    private let context: ModelContext

    init(context: ModelContext) {
        self.context = context
    }

    func save(_ record: RoundRecord) async throws {
        let user = try fetchUserModel(by: record.userIdentifier)
        let model = RoundRecordModel(
            id: record.id,
            userIdentifier: record.userIdentifier,
            sessionIdentifier: record.sessionIdentifier,
            workMinutes: record.workMinutes,
            breakMinutes: record.breakMinutes,
            focusScore: record.focusScore,
            isAborted: record.isAborted,
            createdAt: record.createdAt,
            updatedAt: record.updatedAt,
            user: user
        )
        context.insert(model)
        try context.save()
        print("🧾 RoundRecord saved id=\(model.id) aborted=\(model.isAborted)")
    }

    func fetchRecent(for userIdentifier: String, limit: Int?) async throws -> [RoundRecord] {
        var descriptor = FetchDescriptor<RoundRecordModel>(
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

    private func map(_ model: RoundRecordModel) throws -> RoundRecord {
        return RoundRecord(
            id: model.id,
            userIdentifier: model.userIdentifier,
            sessionIdentifier: model.sessionIdentifier,
            workMinutes: model.workMinutes,
            breakMinutes: model.breakMinutes,
            focusScore: model.focusScore,
            isAborted: model.isAborted,
            createdAt: model.createdAt,
            updatedAt: model.updatedAt
        )
    }
}
