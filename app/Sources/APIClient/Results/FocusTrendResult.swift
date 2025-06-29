//
//  FocusTrendResult.swift
//  SelfPomodoro
//
//  Created by し on 2025/06/15.
//

import Foundation

struct FocusTrendResult: Codable, Identifiable, Equatable {
    var id: Date { date }
    let date: Date
    let focusScore: Int
}
