//
//  StatisticsScreenView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/20.
//

import SwiftUI
import ComposableArchitecture

struct StatisticsScreenView: View {
    let store: StoreOf<ChartFeature>

    var body: some View {
        ChartView(store: store)
    }
}
