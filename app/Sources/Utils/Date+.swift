//
//  Date+.swift
//  SelfPomodoro
//
//  Created by し on 2025/07/03.
//

import Foundation

extension Date {
    static func startOfCurrentWeek() -> Date {
        let calendar = Calendar.current
        let today = calendar.startOfDay(for: Date())
        let weekday = calendar.component(.weekday, from: today)
        let diff = (weekday + 5) % 7
        return calendar.date(byAdding: .day, value: -diff, to: today)!
    }
}
