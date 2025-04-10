//
//  TimerViewModel.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/01/10.
//
import Combine
import Foundation

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
    private let initialTaskDuration: Int // タスクの時間（秒）
    private let initialRestDuration: Int // 休憩の時間（秒）

    init(totalRounds: Int, taskDuration: Int = 10, restDuration: Int = 10) {
        initialTaskDuration = taskDuration
        initialRestDuration = restDuration
        let initialTimer = TimerModel(round: 1, totalRounds: totalRounds, taskDuration: taskDuration, restDuration: restDuration)
        timer = initialTimer
        state = initialTimer.state
        round = initialTimer.round
        self.totalRounds = initialTimer.totalRounds
        totalTaskDuration = TimerModel.formatTime(taskDuration)
        totalRestDuration = TimerModel.formatTime(restDuration)
        let initialTimeRemaining = initialTimer.getCurrentDuration()
        timeRemaining = initialTimeRemaining
        formattedTime = TimerModel.formatTime(initialTimeRemaining)
    }

    func startTimer() {
        isActive = true
        startCountdown()
    }

    func stopTimer() {
        isActive = false
        stopCountdown()
    }

    func setTimer() {
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
            timeRemaining = state == .task ? initialTaskDuration : initialRestDuration
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
        formattedTime = TimerModel.formatTime(timeRemaining)
    }

    private func updatePublishedProperties() {
        state = timer.state
        round = timer.round
        totalRounds = timer.totalRounds
    }
}
