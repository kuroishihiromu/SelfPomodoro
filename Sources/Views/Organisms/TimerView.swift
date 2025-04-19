//
//  TimerView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/01/10.
//

import SwiftUI
import ComposableArchitecture

struct TimerView: View {
    let store: StoreOf<TimerFeature>

    var body: some View {
        WithViewStore(store, observe: { $0 }) { viewStore in
            VStack(spacing: 24) {
                TaskTimerView(
                    currentSeconds: viewStore.currentSeconds,
                    totalSeconds: viewStore.totalSeconds
                )

                HStack(spacing: 20) {
                    if viewStore.isRunning {
                        NormalButton(
                            text: "Stop",
                            bgColor: ColorTheme.skyBlue,
                            fontColor: ColorTheme.white,
                            width:200,
                            height: 52,
                            action: {viewStore.send(.stop)}
                        )
                    } else {
                        NormalButton(
                            text: "Start",
                            bgColor: ColorTheme.navy,
                            fontColor: ColorTheme.white,
                            width:200,
                            height: 52,
                            action: {viewStore.send(.start)}
                        )
                    }
                }
            }
            .padding()
        }
    }
}

#Preview {
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
