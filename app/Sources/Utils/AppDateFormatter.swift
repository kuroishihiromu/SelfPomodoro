//
//  AppDateFormatter.swift
//  SelfPomodoro
//
//  Created by し on 2025/06/15.
//

import Foundation

enum AppDateFormatter {
    static let yyyyMMdd: DateFormatter = {
        let formatter = DateFormatter()
        formatter.dateFormat = "yyyy-MM-dd"
        formatter.locale = Locale(identifier: "en_US_POSIX")
        formatter.timeZone = TimeZone(secondsFromGMT: 0)
        return formatter
    }()
}
