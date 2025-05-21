//
//  TimerScreenView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/01/10.
//

import ComposableArchitecture
import SwiftUI

struct TimerScreenView: View {
    let store: StoreOf<TimerScreenFeature>

    var body: some View {
        WithViewStore(store, observe: \.self ) { _ in
            ZStack {
                VStack(spacing: 20) {
                    TimerView(store: store.scope(state: \.timer, action: \.timer))
                }

                // 評価モーダル
                IfLetStore(store.scope(state: \.evalModal, action: \.evalModal)) { modalStore in
                    EvalModalView(store: modalStore)
                        .background(Color.black.opacity(0.3).ignoresSafeArea())
                }
            }
        }
    }
}
