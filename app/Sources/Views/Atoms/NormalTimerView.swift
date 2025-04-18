//
//  NormalTimerView.swift
//  SelfPomodoro
//
//  Created by し on 2025/03/30.
//

import SwiftUI

struct TaskTimerView: View {
    let currentSeconds: Int
    let totalSeconds: Int
    private let barWidth: CGFloat = 350

    var body: some View {
        VStack(spacing: 12) {
            Text("\(format(currentSeconds)) / \(format(totalSeconds))")
                .font(.system(size: 32, weight: .bold))
                .foregroundColor(.black)

            ZStack(alignment: .leading) {
                // 背景バー
                Capsule()
                    .fill(ColorTheme.lightSkyBlue)
                    .frame(width: barWidth, height: 7)

                // 進捗バー
                Capsule()
                    .fill(ColorTheme.navy)
                    .frame(
                        width: barWidth * progress(current: currentSeconds, total: totalSeconds),
                        height: 7
                    )
            }
        }
    }

    private func format(_ seconds: Int) -> String {
        String(format: "%02d:%02d", seconds / 60, seconds % 60)
    }

    private func progress(current: Int, total: Int) -> Double {
        guard total > 0 else { return 0 }
        return min(Double(current) / Double(total), 1.0)
    }
}

#Preview {
    VStack(spacing: 40) {
        TaskTimerView(currentSeconds: 0, totalSeconds: 60)
        TaskTimerView(currentSeconds: 10, totalSeconds: 60)
        TaskTimerView(currentSeconds: 30, totalSeconds: 60)
        TaskTimerView(currentSeconds: 60, totalSeconds: 60)
    }
    .padding()
}
