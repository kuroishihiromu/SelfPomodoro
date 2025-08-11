//
//  JSONDecoder+Extension.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/05/21.
//

import Foundation

enum AppDecoder {
    static let `default`: JSONDecoder = {
        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase

        // 柔軟な dateDecodingStrategy に変更
        decoder.dateDecodingStrategy = .custom { decoder in
            let container = try decoder.singleValueContainer()
            let dateStr = try container.decode(String.self)
            
            if let date = AppDateFormatter.yearMonthDay.date(from: dateStr) {
                return date
            }

            // ミリ秒あり
            if let date = AppDateFormatter.iso8601WithMillis.date(from: dateStr) {
                return date
            }

            // ミリ秒なし
            if let date = AppDateFormatter.iso8601WithoutMillis.date(from: dateStr) {
                return date
            }

            throw DecodingError.dataCorruptedError(
                in: container,
                debugDescription: "Unrecognized date format: \(dateStr)"
            )
        }

        return decoder
    }()
}
