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
    @ObservableState
    struct State: Equatable {
        var chart = ChartFeature.State()
        var heatMap = HeatMapFeature.State()
    }
    
    enum Action {
        case chart(ChartFeature.Action)
        case heatMap(HeatMapFeature.Action)
    }
    
    var body: some ReducerOf<Self> {
        Scope(state: \.chart, action: \.chart) {
            ChartFeature()
        }
        Scope(state: \.heatMap, action: \.heatMap) {
            HeatMapFeature()
        }
    }
}

struct StatisticsScreenView: View {
    let store: StoreOf<StatisticsFeature>

    var body: some View {
        ScrollView {
            VStack(spacing: 24) {
                ChartView(store: store.scope(state: \.chart, action: \.chart))
                
                HeatMapView(store: store.scope(state: \.heatMap, action: \.heatMap))
            }
            .padding(.bottom)
        }
    }
}
