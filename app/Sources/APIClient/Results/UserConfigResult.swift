//
//  ConfigResult.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/07/03.
//

import Foundation
struct UserConfigResult: Codable, Identifiable, Equatable {
    let id: UUID
    let roundWorkTime: Int
    let roundBreakTime: Int
    let sessionRounds: Int
    let sessionBreakTime: Int
    
    
    enum CodingKeys: String, CodingKey {
        case id = "user_id"
        case roundWorkTime = "round_work_time"
        case roundBreakTime = "round_break_time"
        case sessionRounds = "session_rounds"
        case sessionBreakTime = "session_break_time"
    }
}
