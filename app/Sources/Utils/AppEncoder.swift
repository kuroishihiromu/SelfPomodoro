//
//  AppEncoder.swift
//  SelfPomodoro
//
//  Created by し on 2025/06/15.
//

import Foundation

enum AppEncoder {
    static let `default`: JSONEncoder = {
        let encoder = JSONEncoder()
        encoder.keyEncodingStrategy = .convertToSnakeCase
        encoder.dateEncodingStrategy = .formatted(AppDateFormatter.yyyyMMdd)
        return encoder
    }()
}
