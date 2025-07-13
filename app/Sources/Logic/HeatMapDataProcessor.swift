//
//  HeatMapDataProcessor.swift
//  SelfPomodoro
//
//  Created by し on 2025/07/13.
//

import Foundation

enum HeatMapDataProcessor {
    static func generateCalendarData(for month: Date) -> [[Date?]] {
        let calendar = Calendar.current

        guard let range = calendar.range(of: .day, in: .month, for: month),
              let firstDate = calendar.date(from: calendar.dateComponents([.year, .month], from: month)) else {
            return []
        }

        var days: [Date?] = []

        // 月初の曜日補正（週の何日目から始まるか）
        let firstWeekday = calendar.dateComponents([.weekday], from: firstDate).weekday ?? 1
        let weekdayIndex = (firstWeekday - calendar.firstWeekday + 7) % 7
        days += Array(repeating: nil, count: weekdayIndex)

        // 日付を配列に追加
        for day in range {
            if let baseDate = calendar.date(byAdding: .day, value: day - 1, to: firstDate),
               let date = calendar.date(bySettingHour: 8, minute: 0, second: 0, of: baseDate) {
                days.append(date)
            }
        }


        // 月末が7の倍数で終わらなければ空マスを追加
        let remainder = days.count % 7
        if remainder != 0 {
            days += Array(repeating: nil, count: 7 - remainder)
        }

        // 週ごとに7日単位で分割
        return stride(from: 0, to: days.count, by: 7).map {
            Array(days[$0..<min($0 + 7, days.count)])
        }
    }

    static func dayNumber(for date: Date) -> Int {
        return Calendar.current.component(.day, from: date)
    }
}
