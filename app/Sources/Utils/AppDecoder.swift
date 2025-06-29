//
//  AppDecoder.swift
//  SelfPomodoro
//
//  Created by し on 2025/06/15.
//

import Foundation

enum AppDecoder {
    static let `default`: JSONDecoder = {
        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        decoder.dateDecodingStrategy = .formatted(AppDateFormatter.yyyyMMdd)
        return decoder
    }()
}
