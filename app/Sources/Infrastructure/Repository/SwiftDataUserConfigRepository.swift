//
//  SwiftDataUserConfigRepository.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/07.
//

import Foundation
import SwiftData

@MainActor
final class SwiftDataUserConfigRepository: UserConfigRepository {
    private let context: ModelContext

    init(context: ModelContext) {
        self.context = context
    }

    func fetchLatest(for userIdentifier: String) async throws -> UserConfig? {
        let model = try fetchLatestModel(by: userIdentifier)
        guard let model else { return nil }
        return try mapLatest(model)
    }

    func upsertLatest(_ config: UserConfig) async throws {
        let user = try fetchUserModel(by: config.userIdentifier)
        if let existing = try fetchLatestModel(by: config.userIdentifier) {
            existing.roundWorkMinutes = config.roundWorkMinutes
            existing.roundBreakMinutes = config.roundBreakMinutes
            existing.sessionRounds = config.sessionRounds
            existing.sessionBreakMinutes = config.sessionBreakMinutes
            existing.updatedAt = config.updatedAt
        } else {
            let model = UserConfigLatestModel(
                id: config.id,
                userIdentifier: config.userIdentifier,
                roundWorkMinutes: config.roundWorkMinutes,
                roundBreakMinutes: config.roundBreakMinutes,
                sessionRounds: config.sessionRounds,
                sessionBreakMinutes: config.sessionBreakMinutes,
                updatedAt: config.updatedAt,
                user: user
            )
            context.insert(model)
        }
        try context.save()
        print("🧩 UserConfig latest updated for \(config.userIdentifier)")
    }

    func addRoundHistory(_ entry: UserConfigRoundHistoryEntry) async throws {
        let user = try fetchUserModel(by: entry.userIdentifier)
        let model = UserConfigRoundHistoryModel(
            id: entry.id,
            userIdentifier: entry.userIdentifier,
            recommendedWorkMinutes: entry.recommendedWorkMinutes,
            recommendedBreakMinutes: entry.recommendedBreakMinutes,
            createdAt: entry.createdAt,
            user: user
        )
        context.insert(model)
        try context.save()
        print("🧾 Round history appended for \(entry.userIdentifier)")
    }

    func addSessionHistory(_ entry: UserConfigSessionHistoryEntry) async throws {
        let user = try fetchUserModel(by: entry.userIdentifier)
        let model = UserConfigSessionHistoryModel(
            id: entry.id,
            userIdentifier: entry.userIdentifier,
            recommendedSessionRounds: entry.recommendedSessionRounds,
            recommendedSessionBreakMinutes: entry.recommendedSessionBreakMinutes,
            createdAt: entry.createdAt,
            user: user
        )
        context.insert(model)
        try context.save()
        print("🧾 Session history appended for \(entry.userIdentifier)")
    }

    func fetchRoundHistory(for userIdentifier: String, limit: Int?) async throws -> [UserConfigRoundHistoryEntry] {
        var descriptor = FetchDescriptor<UserConfigRoundHistoryModel>(
            predicate: #Predicate { $0.userIdentifier == userIdentifier },
            sortBy: [SortDescriptor(\.createdAt, order: .reverse)]
        )
        if let limit {
            descriptor.fetchLimit = limit
        }
        return try context.fetch(descriptor).map { try mapRoundHistory($0) }
    }

    func fetchSessionHistory(for userIdentifier: String, limit: Int?) async throws -> [UserConfigSessionHistoryEntry] {
        var descriptor = FetchDescriptor<UserConfigSessionHistoryModel>(
            predicate: #Predicate { $0.userIdentifier == userIdentifier },
            sortBy: [SortDescriptor(\.createdAt, order: .reverse)]
        )
        if let limit {
            descriptor.fetchLimit = limit
        }
        return try context.fetch(descriptor).map { try mapSessionHistory($0) }
    }

    // MARK: - Helpers

    private func fetchLatestModel(by identifier: String) throws -> UserConfigLatestModel? {
        var descriptor = FetchDescriptor<UserConfigLatestModel>(
            predicate: #Predicate { $0.userIdentifier == identifier }
        )
        descriptor.fetchLimit = 1
        return try context.fetch(descriptor).first
    }

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

    private func mapLatest(_ model: UserConfigLatestModel) throws -> UserConfig {
        return UserConfig(
            id: model.id,
            userIdentifier: model.userIdentifier,
            roundWorkMinutes: model.roundWorkMinutes,
            roundBreakMinutes: model.roundBreakMinutes,
            sessionRounds: model.sessionRounds,
            sessionBreakMinutes: model.sessionBreakMinutes,
            updatedAt: model.updatedAt
        )
    }

    private func mapRoundHistory(_ model: UserConfigRoundHistoryModel) throws -> UserConfigRoundHistoryEntry {
        return UserConfigRoundHistoryEntry(
            id: model.id,
            userIdentifier: model.userIdentifier,
            recommendedWorkMinutes: model.recommendedWorkMinutes,
            recommendedBreakMinutes: model.recommendedBreakMinutes,
            createdAt: model.createdAt
        )
    }

    private func mapSessionHistory(_ model: UserConfigSessionHistoryModel) throws -> UserConfigSessionHistoryEntry {
        return UserConfigSessionHistoryEntry(
            id: model.id,
            userIdentifier: model.userIdentifier,
            recommendedSessionRounds: model.recommendedSessionRounds,
            recommendedSessionBreakMinutes: model.recommendedSessionBreakMinutes,
            createdAt: model.createdAt
        )
    }
}
