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
            VStack(spacing: 20) {
                TimerView(store: store.scope(state: \.timer, action: \.timer))
            }
            .onAppear {
                viewStore.send(.onAppear)
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)
            .safeAreaInset(edge: .top, spacing: 0) {
                MenuBarView(title: "Pomodoro")
            }
            // ラウンド設定モーダル
            .fullScreenCover(
                isPresented: viewStore.binding(
                    get: \.roundConfigModalIsPresented,
                    send: TimerScreenFeature.Action.toggleConfigModal
                )
            ) {
                RoundConfigModalView(config: viewStore.userConfig, currentRound: viewStore.timer.round) {
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
            // 評価モーダル
            .overlay {
                if let _ = viewStore.evalModal {
                    ZStack {
                        Color.black.opacity(0.3)
                            .ignoresSafeArea()
                            .onTapGesture {
                                viewStore.send(.evalModal(.cancel))
                            }

                        IfLetStore(
                            store.scope(state: \.evalModal, action: \.evalModal)
                        ) { modalStore in
                            EvalModalView(store: modalStore)
                                .transition(.scale)
                        }
                    }
                }
            }
            .animation(.easeInOut, value: viewStore.evalModal != nil)
            .overlay(alignment: .top) {
                ToastView(state: viewStore.toast)
                    .padding(.top, 40)
            }
        }
    }
}
