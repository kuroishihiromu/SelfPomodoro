//
//  PomodoroTimer.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/05/16.
//

public struct PomodoroTimer: Equatable {
    public enum Phase {
        case task, shortBreak, longBreak
    }

    public var currentSeconds: Int
    public var totalSeconds: Int
    public var isRunning: Bool
    public var round: Int
    public var phase: Phase

    public var taskDuration: Int
    public var shortBreakDuration: Int
    public var longBreakDuration: Int
    public var roundsPerSession: Int

    public init(task: Int, short: Int, long: Int, rps: Int) {
        self.taskDuration = task
        self.shortBreakDuration = short
        self.longBreakDuration = long
        self.roundsPerSession = rps
        self.phase = .task
        self.round = 1
        self.totalSeconds = task
        self.currentSeconds = task
        self.isRunning = false
    }

    public mutating func start() {
        isRunning = true
    }

    public mutating func stop() {
        isRunning = false
    }

    public mutating func tick() -> TimerTickResult {
        guard isRunning else { return .stopped }
        currentSeconds -= 1
        if currentSeconds <= 0 {
            stop()
            return .completed
        }
        return .running(remaining: currentSeconds)
    }

    public mutating func advancePhase() {
        if phase == .task {
            if round % roundsPerSession == 0 {
                phase = .longBreak
            } else {
                phase = .shortBreak
            }
        } else {
            phase = .task
            round += 1
        }

        totalSeconds = currentPhaseDuration
        currentSeconds = totalSeconds
    }

    private var currentPhaseDuration: Int {
        switch phase {
        case .task: return taskDuration
        case .shortBreak: return shortBreakDuration
        case .longBreak: return longBreakDuration
        }
    }
}

public enum TimerTickResult {
    case running(remaining: Int)
    case completed
    case stopped
}
