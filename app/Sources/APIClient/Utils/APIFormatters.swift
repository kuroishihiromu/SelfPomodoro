//
//  APIFormatters.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/05/21.
//

import Foundation

enum APIFormatters {
    
    /// ISO 8601 formatter with fractional seconds and timezone
    static let iso8601WithFractionalSecondsAndTZ: ISO8601DateFormatter = {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [
            .withInternetDateTime,        // yyyy-MM-dd'T'HH:mm:ssZZZZZ
            .withFractionalSeconds,       // .SSS
            .withTimeZone                 // +09:00 etc
        ]
        return formatter
    }()

    static let iso8601Flexible: ISO8601DateFormatter = {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [
            .withInternetDateTime, // "yyyy-MM-dd'T'HH:mm:ssZ" → これで `"2025-07-03T04:38:19Z"` に対応
            .withFractionalSeconds // fractionalSeconds がついててもOK（あれば処理される）
        ]
        return formatter
    }()
    
    

    /// JSONDecoder with custom ISO8601 handling
    static let jsonDecoder: JSONDecoder = {
        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .useDefaultKeys

        decoder.dateDecodingStrategy = .custom { decoder in
            let container = try decoder.singleValueContainer()
            let dateStr = try container.decode(String.self)

            if let date = iso8601Flexible.date(from: dateStr) {
                return date
            }

            // 秒までしかない形式も許容
            if let fallbackDate = APIFormatters.iso8601WithoutFractional.date(from: dateStr) {
                return fallbackDate
            }

            throw DecodingError.dataCorruptedError(
                in: container,
                debugDescription: "Invalid ISO8601 date format: \(dateStr)"
            )
        }

        return decoder
    }()
    
    static let iso8601WithoutFractional: ISO8601DateFormatter = {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime]
        return formatter
    }()

    static let jsonDecoderWithISOEasyVersion: JSONDecoder = {
        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .useDefaultKeys
        decoder.dateDecodingStrategy = .iso8601
        return decoder
    }()
    
}
