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
        WithViewStore(store, observe: { $0 }) { viewStore in
            VStack(spacing: 0) {
                // タブによって画面切り替え
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
                // タブバーを表示
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
