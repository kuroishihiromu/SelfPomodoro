//
//  SessionRecordRepository.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/07.
//

import Foundation

protocol SessionRecordRepository {
    func save(_ record: SessionRecord) async throws
    func fetchRecent(for userIdentifier: String, limit: Int?) async throws -> [SessionRecord]
}
