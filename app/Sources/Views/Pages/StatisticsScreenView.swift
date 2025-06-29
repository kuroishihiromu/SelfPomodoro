//
//  StatisticsScreenView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/20.
//

import SwiftUI

struct StatisticsScreenView: View {
    var body: some View {
        Text(/*@START_MENU_TOKEN@*/"Hello, World!"/*@END_MENU_TOKEN@*/)
    }
}

#Preview {
    StatisticsScreenView()
}

//import ComposableArchitecture
//
//struct StatisticsScreenView: View {
//    let store: StoreOf<StatisticsFeature>
//
//    var body: some View {
//        WithViewStore(store, observe: { $0 }) { viewStore in
//            ChartView(
//                title: "Concentration Score Trend",
//                state: viewStore.chart,
//                sendActionWithAnimation: { action in
//                    withAnimation(.easeInOut) {
//                        _ = viewStore.send(.chart(action))
//                    }
//                }
//            )
//            .task {
//                viewStore.send(.chart(.fetchData))
//            }
//        }
//    }
//}
//
//#Preview {
//    StatisticsScreenView(
//        store: Store(
//            initialState: StatisticsFeature.State(),
//            reducer: { StatisticsFeature() }
//        )
//    )
//}
