//
//  NormalTimerView.swift
//  SelfPomodoro
//
//  Created by し on 2025/03/30.
//

import SwiftUI

struct TaskTimerView: View {
    let minutes: Int
    let seconds: Int
    let totalMinutes: Int
    let totalSeconds: Int
    let fontSize: CGFloat
    let barWidth: CGFloat
    let barHeight: CGFloat
    let progressColor: Color
    let baseColor: Color

    private let timerColor = ColorTheme.black

    private var totalTime: Int {
        totalMinutes * 60 + totalSeconds
    }

    private var remainingTime: Int {
        minutes * 60 + seconds
    }

    private var progress: Double {
        guard totalTime > 0 else { return 0 }
        return 1.0 - (Double(remainingTime) / Double(totalTime))
    }

    private var formattedTime: String {
        String(format: "%02d:%02d", minutes, seconds)
    }

    private var formattedTotalTime: String {
        String(format: "%02d:%02d", totalMinutes, totalSeconds)
    }

    var body: some View {
        VStack(spacing: 12) {
            // 表示: 残り / 合計
            Text("\(formattedTime) / \(formattedTotalTime)")
                .font(.system(size: fontSize, weight: .bold))
                .foregroundColor(timerColor)

            // 進捗バー
            ZStack(alignment: .leading) {
                Capsule()
                    .fill(baseColor)
                    .frame(width: barWidth, height: barHeight)

                Capsule()
                    .fill(progressColor)
                    .frame(width: CGFloat(progress) * barWidth, height: barHeight)
            }
        }
    }
}
#Preview {
    TaskTimerView(
        minutes: 12,
        seconds: 30,
        totalMinutes: 25,
        totalSeconds: 0,
        fontSize: 32,
        barWidth: 350,
        barHeight: 7,
        progressColor: ColorTheme.navy,
        baseColor: ColorTheme.skyBlue
    )
}
