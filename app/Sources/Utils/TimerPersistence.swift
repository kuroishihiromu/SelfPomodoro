//
//  TimerPersistence.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/08/16.
//

import Foundation

struct TimerPersistenceData: Codable {
    let startTime: Date
    let taskDuration: Int
    let shortBreakDuration: Int
    let longBreakDuration: Int
    let roundsPerSession: Int
    let phase: String  // "task", "shortBreak", "longBreak"
    let round: Int
    let isRunning: Bool
    let currentSeconds: Int
}

class TimerPersistence {
    private static let timerDataKey = "timer_persistence_data"

    static func save(_ data: TimerPersistenceData) {
        do {
            let encoded = try JSONEncoder().encode(data)
            UserDefaults.standard.set(encoded, forKey: timerDataKey)
        } catch {
            print("Failed to save timer data: \(error)")
        }
    }

    static func load() -> TimerPersistenceData? {
        guard let data = UserDefaults.standard.data(forKey: timerDataKey) else {
            return nil
        }

        do {
            return try JSONDecoder().decode(TimerPersistenceData.self, from: data)
        } catch {
            print("Failed to load timer data: \(error)")
            return nil
        }
    }

    static func clear() {
        UserDefaults.standard.removeObject(forKey: timerDataKey)
    }
}
