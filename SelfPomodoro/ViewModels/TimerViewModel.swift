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
    
    @Published var isActive: Bool = false
    @Published var isSet: Bool = false
    
    private var cancellable: AnyCancellable?
    private let taskDuration: Int = 10 // タスクの時間（秒）
    private let restDuration: Int = 10 // 休憩の時間（秒）
    
    init(totalRounds: Int) {
        let initialTimer = TimerModel(round: 1, totalRounds: totalRounds)
        self.timer = initialTimer
        self.state = initialTimer.state
        self.round = initialTimer.round
        self.totalRounds = initialTimer.totalRounds
        self.totalTaskDuration = TimerModel.formatTime(taskDuration)
        self.totalRestDuration = TimerModel.formatTime(restDuration)
        self.timeRemaining = taskDuration
        self.formattedTime = TimerModel.formatTime(taskDuration)
    }
    
    func startTimer(){
        isActive = true
        startCountdown()
    }
    
    func stopTimer(){
        isActive = false
        stopCountdown()
    }
    
    func setTimer(){
        isSet = true
        reset()
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
            formattedTime = TimerModel.formatTime(timeRemaining)
        } else {
            nextState()
            timeRemaining = state == .task ? taskDuration : restDuration
            formattedTime = TimerModel.formatTime(timeRemaining)
        }
    }
    
    func nextState() {
        timer.nextState()
        updatePublishedProperties()
    }
    
    func reset() {
        timer.reset()
        updatePublishedProperties()
        formattedTime = TimerModel.formatTime(taskDuration)
    }
    
    private func updatePublishedProperties() {
        state = timer.state
        round = timer.round
        totalRounds = timer.totalRounds
    }
}