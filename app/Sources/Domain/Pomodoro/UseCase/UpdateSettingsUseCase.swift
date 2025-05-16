//
//  UpdateSettingsUseCase.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/05/16.
//

public struct UpdateSettingsUseCaseImpl: UpdateSettingsUseCase {
    public init() {}

    public func execute(task: Int, short: Int, long: Int, rounds: Int) -> PomodoroTimer {
        return PomodoroTimer(task: task, short: short, long: long, rps: rounds)
    }
}
