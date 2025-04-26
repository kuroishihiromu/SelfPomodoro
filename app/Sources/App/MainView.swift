//
//  MainView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/04/10.
//

import ComposableArchitecture
import SwiftUI

struct MainView: View {
    let store: StoreOf<TabButtonFeature>

    var body: some View {
        WithViewStore(store, observe: \.selectedTabIndex) { viewStore in
            VStack(spacing: 0) {
                Group {
                    switch viewStore.state {
                    case 0:
                        TimerScreenView(
                            store: Store(
                                initialState: TimerScreenFeature.State(
                                    timer: TimerFeature.State(
                                        totalSeconds: 10,
                                        taskDuration: 30,
                                        shortBreakDuration: 10,
                                        longBreakDuration: 20,
                                        roundsPerSession:3
                                    ),
                                    todoList: ToDoListFeature.State(),
                                    evalModal: nil
                                ),
                                reducer: { TimerScreenFeature() }
                            )
                        )
                    case 1:
                        TaskManagementScreenView()
                    case 2:
                        StatisticsScreenView()
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
#Preview {
    MainView(store: Store(initialState: TabButtonFeature.State()) {
        TabButtonFeature()
    }
    )
}
