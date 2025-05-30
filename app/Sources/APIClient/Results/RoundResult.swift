//
//  RoundResult.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/05/21.
//

import Foundation
struct RoundResult: Codable, Identifiable, Equatable {
    let id: UUID
    let sessionId: UUID
    let roundOrder: Int
    let startTime: Date
    var endTime: Date?
    var workTime: Int?
    var breakTime: Int?
    var focusScore: Int?
    var isAborted: Bool

    enum CodingKeys: String, CodingKey {
        case id
        case sessionId = "session_id"
        case roundOrder = "round_order"
        case startTime = "start_time"
        case endTime = "end_time"
        case workTime = "work_time"
        case breakTime = "break_time"
        case focusScore = "focus_score"
        case isAborted = "is_aborted"
    }
}
