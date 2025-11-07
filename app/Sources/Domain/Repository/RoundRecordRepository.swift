//
//  RoundRecordRepository.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/07.
//

import Foundation

protocol RoundRecordRepository {
    func save(_ record: RoundRecord) async throws
    func fetchRecent(for userIdentifier: String, limit: Int?) async throws -> [RoundRecord]
}
