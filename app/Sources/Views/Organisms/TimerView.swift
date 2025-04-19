//
//  TimerView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/01/10.
//

import ComposableArchitecture
import SwiftUI

struct TimerView: View {
    let store: StoreOf<TimerFeature>

    var body: some View {
        WithViewStore(store, observe: { $0 }) { viewStore in
            VStack(spacing: 24) {
                // ラウンド数の表示
                Text("Round \(viewStore.round)")
                    .font(.system(size: 22, weight: .bold))

                // フェーズ（Task/Break）の表示
                Text(viewStore.phase == .task ? "Task Time" : "Break Time")
                    .font(.system(size: 18, weight: .medium))
                    .foregroundColor(.gray)

                // タイマーバー（色切替付き）
                TaskTimerView(
                    currentSeconds: viewStore.currentSeconds,
                    totalSeconds: viewStore.totalSeconds,
                    mainColor: viewStore.phase == .task ? ColorTheme.navy : ColorTheme.skyBlue,
                    subColor: viewStore.phase == .task ? ColorTheme.lightSkyBlue : ColorTheme.navy
                )

                // ボタン（Start / Stop）
                HStack(spacing: 20) {
                    if viewStore.isRunning {
                        NormalButton(
                            text: "Stop",
                            bgColor: ColorTheme.skyBlue,
                            fontColor: ColorTheme.white,
                            width: 200,
                            height: 52,
                            action: { viewStore.send(.stop) }
                        )
                    } else {
                        NormalButton(
                            text: "Start",
                            bgColor: ColorTheme.navy,
                            fontColor: ColorTheme.white,
                            width: 200,
                            height: 52,
                            action: { viewStore.send(.start) }
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
                totalSeconds: 10,
                taskDuration: 30,
                restDuration: 10
            ),
            reducer: { TimerFeature() }
        )
    )
}
