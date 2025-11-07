//
//  User.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/05.
//

import Foundation

struct User: Identifiable, Equatable {
    let id: UUID
    let identifier: String
    var createdAt: Date
    var updatedAt: Date
}
