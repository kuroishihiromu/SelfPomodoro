//
//  DateFormatter+Extensioin.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/05/21.
//

import Foundation

enum AppDateFormatter {
    /// ISO8601形式（ミリ秒あり、タイムゾーンあり）
    static let iso8601WithMillis: DateFormatter = {
        let formatter = DateFormatter()
        formatter.locale = Locale(identifier: "en_US_POSIX")
        formatter.timeZone = TimeZone(secondsFromGMT: 0)
        formatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ss.SSSXXXXX"
        return formatter
    }()
    
    /// 他の形式があればここに追加
}
