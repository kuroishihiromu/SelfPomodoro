//
//  PomodoroTimerUseCase.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/05/16.
//

public struct StartTimerUseCaseImpl: StartTimerUseCase {
    public init() {}
    
    public func execute(timer: PomodoroTimer) -> PomodoroTimer {
        var updated = timer
        updated.start()
        return updated
    }
}
