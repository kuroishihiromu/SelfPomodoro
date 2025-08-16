//
//  LocalizationHelper.swift
//  SelfPomodoro
//
//  Created by Claude on 2025/01/16.
//

import Foundation

/// ローカライゼーション用ヘルパー
extension String {
    /// ローカライズされた文字列を取得
    var localized: String {
        return NSLocalizedString(self, comment: "")
    }
    
    /// パラメータ付きローカライズ
    func localized(with arguments: CVarArg...) -> String {
        return String(format: self.localized, arguments: arguments)
    }
}

/// ローカライゼーション用の列挙型
enum L10n {
    // MARK: - 認証画面
    enum Auth {
        static let createAccount = "auth.create_account".localized
        static let description = "auth.description".localized
        static let emailPlaceholder = "auth.email_placeholder".localized
        static let passwordPlaceholder = "auth.password_placeholder".localized
        static let termsAgreement = "auth.terms_agreement".localized
        static let loginButton = "auth.login_button".localized
        static let signupButton = "auth.signup_button".localized
        static let termsRequired = "auth.terms_required".localized
    }
    
    // MARK: - タイマー画面
    enum Timer {
        static let pomodoro = "timer.pomodoro".localized
        static let shortBreak = "timer.short_break".localized
        static let longBreak = "timer.long_break".localized
        static let start = "timer.start".localized
        static let pause = "timer.pause".localized
        static let resume = "timer.resume".localized
        static let stop = "timer.stop".localized
        static let reset = "timer.reset".localized
    }
    
    // MARK: - タスク管理
    enum Tasks {
        static let title = "tasks.title".localized
        static let addTask = "tasks.add_task".localized
        static let taskPlaceholder = "tasks.task_placeholder".localized
        static let completed = "tasks.completed".localized
        static let pending = "tasks.pending".localized
        static let delete = "tasks.delete".localized
    }
    
    // MARK: - 統計画面
    enum Stats {
        static let title = "stats.title".localized
        static let focusTrend = "stats.focus_trend".localized
        static let dailySummary = "stats.daily_summary".localized
        static let weeklySummary = "stats.weekly_summary".localized
        static let concentrationScore = "stats.concentration_score".localized
    }
    
    // MARK: - プロフィール画面
    enum Profile {
        static let title = "profile.title".localized
        static let settings = "profile.settings".localized
        static let logout = "profile.logout".localized
        static let account = "profile.account".localized
    }
    
    // MARK: - 設定画面
    enum Settings {
        static let title = "settings.title".localized
        static let timerSettings = "settings.timer_settings".localized
        static let workDuration = "settings.work_duration".localized
        static let shortBreakDuration = "settings.short_break_duration".localized
        static let longBreakDuration = "settings.long_break_duration".localized
        static let roundsPerSession = "settings.rounds_per_session".localized
        static let notifications = "settings.notifications".localized
        static let sound = "settings.sound".localized
    }
    
    // MARK: - 一般
    enum General {
        static let cancel = "general.cancel".localized
        static let ok = "general.ok".localized
        static let save = "general.save".localized
        static let delete = "general.delete".localized
        static let edit = "general.edit".localized
        static let done = "general.done".localized
        static let close = "general.close".localized
        static let yes = "general.yes".localized
        static let no = "general.no".localized
    }
    
    // MARK: - エラーメッセージ
    enum Error {
        static let network = "error.network".localized
        static let authentication = "error.authentication".localized
        static let invalidEmail = "error.invalid_email".localized
        static let invalidPassword = "error.invalid_password".localized
        static let unknown = "error.unknown".localized
    }
    
    // MARK: - 成功メッセージ
    enum Success {
        static let login = "success.login".localized
        static let signup = "success.signup".localized
        static let taskAdded = "success.task_added".localized
        static let taskCompleted = "success.task_completed".localized
        static let settingsSaved = "success.settings_saved".localized
    }
}