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
        WithViewStore(store, observe: \.self) { viewStore in
            ZStack {
                VStack(spacing: 20) {
                    TimerView(store: store.scope(state: \.timer, action: \.timer))
                }
                .onAppear {
                    viewStore.send(.onAppear)
                }
                
                // 評価モーダル
                IfLetStore(store.scope(state: \.evalModal, action: \.evalModal)) { modalStore in
                    EvalModalView(store: modalStore)
                        .background(Color.black.opacity(0.3).ignoresSafeArea())
                }
            }
            // ラウンド設定モーダル
            .fullScreenCover(
                isPresented: viewStore.binding(
                    get: \.roundConfigModalIsPresented,
                    send: TimerScreenFeature.Action.toggleConfigModal
                )
            ) {
                RoundConfigModalView(config: viewStore.userConfig, currentRound: store.timer.round) {
                    viewStore.send(.StartRoundButtonTapped)
                    viewStore.send(.toggleConfigModal(false))
                }
            }
            .fullScreenCover(
                isPresented: viewStore.binding(
                    get: \.sessionCompleteModal,
                    send: TimerScreenFeature.Action.toggleSessionCompleteModal
                )
            ) {
                SessionCompleteModalView() {
                    viewStore.send(.sessionCompleteModalTapped)
                    viewStore.send(.toggleSessionCompleteModal(false))
                }
            }

        }
    }
}
