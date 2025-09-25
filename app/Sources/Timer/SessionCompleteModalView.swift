//
//  SessionCompleteModalView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/07/03.
//

import SwiftUI

struct SessionCompleteModalView: View {
    let onStart: () -> Void

    var body: some View {
        VStack(spacing: 20) {
            Text("セッションが完了しました。お疲れ様です。")

            NormalButton(
                text: "閉じる",
                bgColor: ColorTheme.navy,
                fontColor: ColorTheme.white,
                width: 200,
                height: 44,
                action: onStart
            )
        }
        .padding()
        .background(Color.white)
        .cornerRadius(20)
        .shadow(radius: 10)
    }
}
