//
//  StatisticsScreenView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/20.
//

import SwiftUI
import ComposableArchitecture

@Reducer
struct StatisticsFeature {
    enum DisplayMode: Equatable {
        case chart
        case heatMap
    }

    @ObservableState
    struct State: Equatable {
        var chart = ChartFeature.State()
        var heatMap = HeatMapFeature.State()
        var displayMode: DisplayMode = .chart
    }
    
    enum Action {
        case chart(ChartFeature.Action)
        case heatMap(HeatMapFeature.Action)
        case changeMode(DisplayMode)
    }
    
    var body: some ReducerOf<Self> {
        Scope(state: \.chart, action: \.chart) {
            ChartFeature()
        }
        Scope(state: \.heatMap, action: \.heatMap) {
            HeatMapFeature()
        }
        Reduce { state, action in
            switch action {
            case .changeMode(let mode):
                state.displayMode = mode
                return .none
            case .chart, .heatMap:
                return .none
            }
        }
    }
}

struct StatisticsScreenView: View {
    let store: StoreOf<StatisticsFeature>

    var body: some View {
        WithViewStore(store, observe: \.displayMode) { viewStore in
            VStack {
                // 切り替えボタン
                SegmentedTabView(
                    selection: viewStore.binding(get: { $0 }, send: { .changeMode($0) })
                )
                .padding(.horizontal)
                .padding(.top, 12)

                // チャートかヒートマップ
                VStack {
                    if viewStore.state == .chart {
                        ChartView(store: store.scope(state: \.chart, action: \.chart))
                            .transition(.opacity)
                    } else {
                        HeatMapView(store: store.scope(state: \.heatMap, action: \.heatMap))
                            .transition(.opacity)
                    }
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .center)
                .animation(.easeInOut, value: viewStore.state)
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)
            .safeAreaInset(edge: .top, spacing: 0) {
                MenuBarView(title: "Statistics")
            }
        }
    }
}
