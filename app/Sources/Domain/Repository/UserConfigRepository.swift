//
//  UserConfigRepository.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/07.
//

import Foundation

protocol UserConfigRepository {
    func fetchLatest(for userIdentifier: String) async throws -> UserConfig?
    func upsertLatest(_ config: UserConfig) async throws

    func addRoundHistory(_ entry: UserConfigRoundHistoryEntry) async throws
    func addSessionHistory(_ entry: UserConfigSessionHistoryEntry) async throws

    func fetchRoundHistory(for userIdentifier: String, limit: Int?) async throws -> [UserConfigRoundHistoryEntry]
    func fetchSessionHistory(for userIdentifier: String, limit: Int?) async throws -> [UserConfigSessionHistoryEntry]
}
