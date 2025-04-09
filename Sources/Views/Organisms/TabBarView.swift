//
//  TabBarView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/20.
//

import SwiftUI
import ComposableArchitecture

struct TabBarView: View {
    @State private var selectedIndex = 0
    let store: StoreOf<TabButtonFeature>

    var body: some View {
        WithViewStore(store, observe: {$0}) { viewStore in
            VStack(spacing: 0) {
                // 選択中の画面を表示
                Group {
                    switch viewStore.selectedTabIndex {
                    case 0: HomeScreenView()
                    case 1: TaskManagementScreenView()
                    case 2: StatisticsScreenView()
                    case 3: ProfileScreenView()
                    default: EmptyView()
                    }
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity)
                
                // タブバー本体
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
}

#Preview {
    TabBarView(
        store: Store(initialState: TabButtonFeature.State()) {
            TabButtonFeature()
        }
    )
}
