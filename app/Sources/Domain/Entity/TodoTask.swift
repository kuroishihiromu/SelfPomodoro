//
//  Task.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/05.
//

import Foundation

struct TodoTask: Identifiable, Equatable {
    let id: UUID
    var detail: String
    var isCompleted: Bool
    var createdAt: Date
    var updatedAt: Date
}
