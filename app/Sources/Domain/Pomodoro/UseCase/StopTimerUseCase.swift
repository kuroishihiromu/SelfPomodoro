//
//  StopTimerUseCase.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/05/16.
//

public struct StopTimerUseCaseImpl: StopTimerUseCase {
    public init() {}
    
    public func execute(timer: PomodoroTimer) -> PomodoroTimer {
        var updated = timer
        updated.stop()
        return updated
    }
}
