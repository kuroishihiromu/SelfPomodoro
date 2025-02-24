//
//  TimerModel.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/01/10.
//

import Foundation

struct TimerModel {
    enum TimerState {
        case task
        case rest
    }
    
    var state: TimerState = .task  // 初期状態を Task に設定
    var round: Int     // 現在のラウンド数
    var totalRounds: Int  // 最大ラウンド数
    var taskDuration: Int
    var restDuration: Int
    
    mutating func nextState() {
        // タイマーの状態を切り替え
        if state == .task {
            state = .rest
        } else {
            state = .task
            round += 1
        }
        
        // ラウンドが最大数に達した場合の処理
        if round > totalRounds {
            reset()
        }
    }
    
    func getCurrentDuration() -> Int {
        return state == .task ? taskDuration : restDuration
    }
    
    mutating func reset() {
        state = .task
        round = 1
    }
    
    // 時間をフォーマットする関数
    static func formatTime(_ seconds: Int) -> String {
        let minutes = seconds / 60
        let remainingSeconds = seconds % 60
        return String(format: "%02d:%02d", minutes, remainingSeconds)
    }
}
