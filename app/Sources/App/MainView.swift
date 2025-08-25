//
//  MainView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/04/10.
//

import ComposableArchitecture
import SwiftUI

struct MainView: View {
    var token: AuthTokens
    let authStore: StoreOf<AuthFeature>
    
    let store: StoreOf<TabButtonFeature> = Store(initialState: TabButtonFeature.State()) {
        TabButtonFeature()
    }
    
    init(token: AuthTokens, authStore: StoreOf<AuthFeature>) {
        self.token = token
        self.authStore = authStore
    }
    
    let timerStore = Store(
        initialState: TimerScreenFeature.State(
            timer: TimerFeature.State(
                totalSeconds: 23*62,
                taskDuration: 30,
                shortBreakDuration: 5*60,
                longBreakDuration: 20,
                roundsPerSession: 3
            ),
            evalModal: nil
        ),
        reducer: { TimerScreenFeature() }
    )
    let toDoStore = Store(initialState: ToDoListFeature.State()) {
        ToDoListFeature()
    }
    
    let statisticsStore = Store(initialState: StatisticsFeature.State()) {
        StatisticsFeature()
    }
    
    var body: some View {
        WithViewStore(store, observe: \.selectedTabIndex) { viewStore in
            VStack(spacing: 0) {
                switch viewStore.state {
                case 0:
                    TimerScreenView(store: timerStore)
                        .frame(maxWidth: .infinity, maxHeight: .infinity)
                case 1:
                    TaskManagementScreenView(
                        store: store.scope(state: \.todoListState, action: \.todoList)
                    )
                    .frame(maxWidth: .infinity, maxHeight: .infinity)
                case 2:
                    StatisticsScreenView(store: statisticsStore)
                        .frame(maxWidth: .infinity, maxHeight: .infinity)
                case 3:
                    ProfileScreenView(authStore: authStore)
                default:
                    EmptyView()
                }
                TabBarView(store: store)
            }
            .onReceive(NotificationCenter.default.publisher(for: UIApplication.willTerminateNotification)) { _ in
                // アプリ完全終了時のみ永続化
                timerStore.send(.timer(.saveTimerState))
            }
        }
    }
}
