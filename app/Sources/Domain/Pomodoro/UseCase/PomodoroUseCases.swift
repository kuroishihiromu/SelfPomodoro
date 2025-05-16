//
//  PomodoroUseCases.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/05/16.
//

//Timerをスタート
public protocol StartTimerUseCase {
    func execute(timer: PomodoroTimer) -> PomodoroTimer
}

//Timerをストップ
public protocol StopTimerUseCase {
    func execute(timer: PomodoroTimer) -> PomodoroTimer
}

//タイマーを進める
public protocol TickTimerUseCase {
    func execute(timer: PomodoroTimer) -> (PomodoroTimer, TimerTickResult)
}

//タイマーの設定を更新
public protocol UpdateSettingsUseCase {
    func execute(task: Int, short: Int, long: Int, rounds: Int) -> PomodoroTimer
}

// Phaseが完了
public protocol AdvancePhaseUseCase {
    func execute(timer: PomodoroTimer) -> PomodoroTimer
}
