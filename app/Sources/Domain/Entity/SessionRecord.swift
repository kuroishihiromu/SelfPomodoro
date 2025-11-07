//
//  SessionRecord.swift
//  SelfPomodoro
//
//  Created by Codex on 2025/02/15.
//

import Foundation

struct SessionRecord: Identifiable, Equatable {
    let id: UUID
    let userIdentifier: String
    let sessionIdentifier: UUID
    var roundCount: Int
    var totalWorkMinutes: Double
    var breakMinutes: Double
    var averageFocusScore: Double
    var isAborted: Bool
    var createdAt: Date
    var updatedAt: Date
}
