//
//  DebugRoundRecordRepository.swift
//  SelfPomodoro
//
//  Created by Codex on 2025/11/07.
//

#if DEBUG
import Foundation

@MainActor
final class DebugRoundRecordRepository: RoundRecordRepository {
    private let base: any RoundRecordRepository
    private let fixedCount: Int
    private let fixedTotalMinutes: Double

    init(base: any RoundRecordRepository, fixedCount: Int = 80, fixedTotalMinutes: Double = 1200) {
        self.base = base
        self.fixedCount = fixedCount
        self.fixedTotalMinutes = fixedTotalMinutes
    }

    func save(_ record: RoundRecord) async throws {
        try await base.save(record)
    }

    func fetchRecent(for userIdentifier: String, limit: Int?) async throws -> [RoundRecord] {
        try await base.fetchRecent(for: userIdentifier, limit: limit)
    }

    func countCompleted(for userIdentifier: String) async throws -> Int {
        fixedCount
    }

    func totalWorkMinutes(for userIdentifier: String) async throws -> Double {
        fixedTotalMinutes
    }
}
#endif
