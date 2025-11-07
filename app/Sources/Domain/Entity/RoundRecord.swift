//
//  RoundRecord.swift
//  SelfPomodoro
//
//  Created by Codex on 2025/02/15.
//

import Foundation

struct RoundRecord: Identifiable, Equatable {
    let id: UUID
    let userIdentifier: String
    let sessionIdentifier: UUID
    var workMinutes: Double
    var breakMinutes: Double
    var focusScore: Int?
    var isAborted: Bool
    var createdAt: Date
    var updatedAt: Date
}
