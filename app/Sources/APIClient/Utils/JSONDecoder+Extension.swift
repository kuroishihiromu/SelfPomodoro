//
//  JSONDecoder+Extension.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/05/21.
//

import Foundation

enum AppDecoder {
    /// API用の共通デコーダー
    static let `default`: JSONDecoder = {
        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        decoder.dateDecodingStrategy = .formatted(AppDateFormatter.iso8601WithMillis)
        return decoder
    }()
}
