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

    /// JSONDecoder with custom ISO8601 handling
    static let jsonDecoder: JSONDecoder = {
        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .useDefaultKeys

        decoder.dateDecodingStrategy = .custom { decoder in
            let container = try decoder.singleValueContainer()
            let dateStr = try container.decode(String.self)

            guard let date = iso8601WithFractionalSecondsAndTZ.date(from: dateStr) else {
                throw DecodingError.dataCorruptedError(
                    in: container,
                    debugDescription: "Invalid ISO8601 date format: \(dateStr)"
                )
            }

            return date
        }

        return decoder
    }()
}
