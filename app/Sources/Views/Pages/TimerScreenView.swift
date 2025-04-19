//
//  TimerScreenView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/01/10.
//

import ComposableArchitecture
import SwiftUI

struct TimerScreenView: View {

    var body: some View {
        ZStack{
            VStack{
                TimerView(
                    store: Store(
                        initialState: TimerFeature.State(
                            totalSeconds: 30,
                            taskDuration: 30,
                            restDuration: 10
                        ),
                        reducer: { TimerFeature() }
                    )
                )
                ToDoListView(store: Store(
                    initialState: ToDoListFeature.State(),
                    reducer: {
                        ToDoListFeature()
                    }
                ))
            }
        }
    }
}

#Preview {
    TimerScreenView()
}
