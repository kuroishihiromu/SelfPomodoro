//
//  FocusData.swift
//  SelfPomodoro
//
//  Created by し on 2025/07/13.
//

import Foundation

struct FocusData: Equatable, Identifiable, Decodable {
    let id = UUID()
    let date: Date
    let hour: Int
    let focus_score: Int
}

struct FocusDataWrapper: Decodable {
    let items: [FocusData]
}
