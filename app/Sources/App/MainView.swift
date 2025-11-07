//
//  MainView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/04/10.
//

import ComposableArchitecture
import SwiftUI

struct MainView: View {
    private let store: StoreOf<TabButtonFeature>
    private let timerStore: StoreOf<TimerScreenFeature>
    private let toDoStore: StoreOf<ToDoListFeature>
    private let statisticsStore: StoreOf<StatisticsFeature>

    init(dependencies: DependencyValues) {
        store = withDependencies {
            $0 = dependencies
        } operation: {
            Store(initialState: TabButtonFeature.State()) {
                TabButtonFeature()
            }
        }

        timerStore = withDependencies {
            $0 = dependencies
        } operation: {
            Store(
                initialState: TimerScreenFeature.State(
                    timer: TimerFeature.State(
                        totalSeconds: 23 * 62,
                        taskDuration: 30,
                        shortBreakDuration: 5 * 60,
                        longBreakDuration: 20,
                        roundsPerSession: 3
                    ),
                    evalModal: nil
                )
            ) {
                TimerScreenFeature()
            }
        }

        toDoStore = withDependencies {
            $0 = dependencies
        } operation: {
            Store(initialState: ToDoListFeature.State()) {
                ToDoListFeature()
            }
        }

        statisticsStore = withDependencies {
            $0 = dependencies
        } operation: {
            Store(initialState: StatisticsFeature.State()) {
                StatisticsFeature()
            }
        }
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
                    ProfileScreenView()
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
