//
//  TimerViewModel.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/01/10.
//
import Foundation
import Combine

class TimerViewModel: ObservableObject {
    @Published private var timer: TimerModel
    @Published var state: TimerModel.TimerState
    @Published var round: Int
    @Published var totalRounds: Int
    @Published var timeRemaining: Int
    @Published var formattedTime: String
    @Published var totalTaskDuration: String
    @Published var totalRestDuration: String
    
    private var cancellable: AnyCancellable?
    private let taskDuration: Int = 25*60 // タスクの時間（秒）
    private let restDuration: Int = 5*60 // 休憩の時間（秒）
    
    init(totalRounds: Int) {
        let initialTimer = TimerModel(round: 1, totalRounds: totalRounds)
        self.timer = initialTimer
        self.state = initialTimer.state
        self.round = initialTimer.round
        self.totalRounds = initialTimer.totalRounds
        self.totalTaskDuration = Self.formatTime(taskDuration)
        self.totalRestDuration = Self.formatTime(restDuration)
        self.timeRemaining = taskDuration
        self.formattedTime = Self.formatTime(taskDuration)
    }
    
    func startCountdown() {
        cancellable = Timer
            .publish(every: 1, on: .main, in: .common)
            .autoconnect()
            .sink { [weak self] _ in
                self?.updateTimer()
            }
    }
    
    func stopCountdown() {
        cancellable?.cancel()
    }
    
    private func updateTimer() {
        if timeRemaining > 0 {
            timeRemaining -= 1
            formattedTime = Self.formatTime(timeRemaining)
        } else {
            nextState()
            timeRemaining = state == .task ? taskDuration : restDuration
            formattedTime = Self.formatTime(timeRemaining)
        }
    }
    
    func nextState() {
        timer.nextState()
        updatePublishedProperties()
    }
    
    func reset() {
        timer.reset()
        updatePublishedProperties()
        formattedTime = Self.formatTime(taskDuration)
    }
    
    private static func formatTime(_ seconds: Int) -> String {
        let minutes = seconds / 60
        let remainingSeconds = seconds % 60
        return String(format: "%02d:%02d", minutes, remainingSeconds)
    }
    
    private func updatePublishedProperties() {
        state = timer.state
        round = timer.round
        totalRounds = timer.totalRounds
    }
}
