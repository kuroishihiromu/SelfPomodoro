//
//  RoundConfigModalView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/07/03.
//

import SwiftUI

struct RoundConfigModalView: View {
    let config: UserConfigResult
    let currentRound: Int
    let onStart: () -> Void

    var body: some View {
        VStack(spacing: 20) {
            Text("作業時間: \(config.roundWorkTime)分")
            Text("休憩時間: \(config.roundBreakTime)分")
            Text(" \(currentRound) / \(config.sessionRounds)")

            NormalButton(
                text: "Start Round",
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
        .onAppear {
            print("🪟 RoundConfigModalView show: work=\(config.roundWorkTime), break=\(config.roundBreakTime), currentRound=\(currentRound)/\(config.sessionRounds)")
        }
    }
}
