//
//  StatisticsScreenPreview.swift
//  SelfPomodoro
//
//  Created by し on 2025/07/03.
//

import SwiftUI
import ComposableArchitecture

struct StatisticsScreenPreview: View {
    let store: StoreOf<StatisticsFeature>

    var body: some View {
        WithViewStore(store, observe: { $0 }) { viewStore in
            ChartView(
                title: "Concentration Score Trend",
                state: viewStore.chart,
                sendActionWithAnimation: { action in
                    withAnimation(.easeInOut) {
                        _ = viewStore.send(.chart(action))
                    }
                }
            )
            .task {
                viewStore.send(.chart(.fetchFocusTrend))
            }
        }
    }
}

#Preview {
    StatisticsScreenPreview(
        store: Store(
            initialState: StatisticsFeature.State(),
            reducer: { StatisticsFeature() }
        )
    )
}
