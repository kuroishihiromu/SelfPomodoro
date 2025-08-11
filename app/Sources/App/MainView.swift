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
    
    let store: StoreOf<TabButtonFeature> = Store(initialState: TabButtonFeature.State()) {
        TabButtonFeature()
    }
    
    init(token: AuthTokens) {
        self.token = token
        print("token: \(self.token)")
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
                Group {
                    switch viewStore.state {
                    case 0:
                        TimerScreenView(store: timerStore)
                    case 1:
                        TaskManagementScreenView(
                            store: store.scope(state: \.todoListState, action: \.todoList)
                        )
                    case 2:
                        StatisticsScreenView(store: statisticsStore)

                    case 3:
                        ProfileScreenView()
                    default:
                        EmptyView()
                    }
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity)

                TabBarView(store: store)
            }
        }
    }
}
