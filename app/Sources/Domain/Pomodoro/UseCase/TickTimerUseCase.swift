//
//  TickTimerUseCase.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/05/16.
//

public struct TickTimerUseCaseImpl: TickTimerUseCase {
    public init() {}
    
    public func execute(timer: PomodoroTimer) -> (PomodoroTimer, TimerTickResult) {
        var updated = timer
        let result = updated.tick()
        return (updated, result)
    }
}
