//
//  TimerScreenView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/01/10.
//

import ComposableArchitecture
import SwiftUI

struct TimerScreenView: View {
//    var timerViewModel: TimerViewModel

    var body: some View {
        TimerView(
            store: Store(
                initialState: TimerFeature.State(
                    currentSeconds: 0,
                    totalSeconds: 10
                ),
                reducer: { TimerFeature() }
            )
        )
    }
}

#Preview {
    TabBarView(
        store: Store(initialState: TabButtonFeature.State()) {
            TabButtonFeature()
        }
    )
}
