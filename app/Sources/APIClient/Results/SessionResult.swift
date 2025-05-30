//
//  SessionResult.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/05/21.
//

import Foundation
struct SessionResult: Codable, Identifiable, Equatable {
    let id: UUID
    let startTime: Date
    var endTime: Date?
    var averageFocus: Double?
    var totalWorkMin: Int?
    var roundCount: Int?
    var breakTime: Int?

    enum CodingKeys: String, CodingKey {
        case id
        case startTime = "start_time"
        case endTime = "end_time"
        case averageFocus = "average_focus"
        case totalWorkMin = "total_work_min"
        case roundCount = "round_count"
        case breakTime = "break_time"
    }
}
