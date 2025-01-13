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
    
    mutating func reset() {
        state = .task
        round = 1
    }
}
