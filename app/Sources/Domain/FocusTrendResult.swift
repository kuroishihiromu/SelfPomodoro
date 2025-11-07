//
//  FocusTrendResult.swift
//  SelfPomodoro
//
//  Created by し on 2025/06/15.
//

import Foundation
public struct FocusTrendResponse: Decodable {
    let items: [FocusTrendResult]
}

struct FocusTrendResult: Equatable, Decodable {
    let date: Date
    let focusScore: Double

    enum CodingKeys: String, CodingKey {
        case date
        case focusScore
    }
}
