//
//  TabBarView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/20.
//

import ComposableArchitecture
import SwiftUI

struct TabBarView: View {
    let store: StoreOf<TabButtonFeature>

    var body: some View {
        WithViewStore(store, observe: { $0 }) { viewStore in
            HStack {
                TabButton(
                    icon: "gauge.with.needle",
                    isSelected: viewStore.selectedTabIndex == 0,
                    action: { viewStore.send(.homeButtonTapped) }
                )
                TabButton(
                    icon: "checkmark.circle",
                    isSelected: viewStore.selectedTabIndex == 1,
                    action: { viewStore.send(.tasksButtonTapped) }
                )
                TabButton(
                    icon: "chart.bar",
                    isSelected: viewStore.selectedTabIndex == 2,
                    action: { viewStore.send(.statsButtonTapped) }
                )
                TabButton(
                    icon: "person.crop.circle",
                    isSelected: viewStore.selectedTabIndex == 3,
                    action: { viewStore.send(.profileButtonTapped) }
                )
            }
            .padding()
            .background(.ultraThinMaterial)
        }
    }
}

#Preview {
    TabBarView(
        store: Store(initialState: TabButtonFeature.State()) {
            TabButtonFeature()
        }
    )
}
