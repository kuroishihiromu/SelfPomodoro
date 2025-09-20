//
//  EvalModalView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/20.
//

import ComposableArchitecture
import SwiftUI

struct EvalModalView: View {
    let store: StoreOf<EvalModalFeature>

    var body: some View {
        WithViewStore(store, observe: { $0 }) { viewStore in
            ZStack {
                VStack(spacing: 24) {
                    // タイトル
                    Text("このラウンドはいかがでしたか？")
                        .font(.system(size: 25, weight: .bold))
                        .padding(.bottom, 16)

                    Text("集中レベルを入力してください。")
                        .font(.system(size: 12, weight: .semibold))
                        .padding(.bottom, 10)

                    // ラウンド情報
                    Text("\("ラウンド") \(viewStore.round)")
                        .font(.system(size: 22, weight: .bold))

                    // スコアバー
                    VStack(spacing: 12) {
                        Text("\(Int(viewStore.score * 100)) / 100")
                            .font(.system(size: 30, weight: .bold))

                        EvalBar(
                            value: viewStore.binding(
                                get: \.score,
                                send: EvalModalFeature.Action.updateScore
                            ),
                            width: 200,
                            height: 20
                        )
                    }

                    // 開始ボタン
                    NormalButton(
                        text: "休憩開始",
                        bgColor: ColorTheme.navy,
                        fontColor: ColorTheme.white,
                        width: 180,
                        height: 48,
                        action: {
                            viewStore.send(.submitEval(viewStore.score))
                        }
                    )
                }
                .padding(EdgeInsets(top: 25, leading: 70, bottom: 35, trailing: 70))
                .overlay(
                    RoundedRectangle(cornerRadius: 3)
                        .stroke(ColorTheme.navy, lineWidth: 4)
                )
                .background(ColorTheme.white)
                .frame(width: 650, height: 650)
            }
        }
    }
}
