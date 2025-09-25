//
//  PieChartView.swift
//  SelfPomodoro
//
//  Created by し on 2025/04/19.
//

import SwiftUI

struct PieChartView: View {
    let completed: Double

    var body: some View {
        VStack(spacing: 16) {
            // 円グラフ
            GeometryReader { geometry in
                let clampedCompleted = max(0, min(completed, 100)) // 0~100の範囲に収まるように処理
                let completedAngle = 360 * clampedCompleted / 100

                ZStack {
                    // 達成済み
                    PieSlice(
                        startAngle: .degrees(-90),
                        endAngle: .degrees(completedAngle - 90)
                    )
                    .fill(ColorTheme.navy)
                    // 未達成
                    PieSlice(
                        startAngle: .degrees(completedAngle - 90),
                        endAngle: .degrees(270)
                    )
                    .fill(ColorTheme.darkGray)
                }
            }
            .aspectRatio(1, contentMode: .fit)
            .frame(width: 200, height: 200)

            // ラベル
            HStack(spacing: 24) {
                HStack(spacing: 8) {
                    Circle()
                        .fill(ColorTheme.navy).frame(width: 12, height: 12)
                    Text("達成済み")
                        .font(.system(size: 14)).foregroundColor(ColorTheme.black)
                }
                HStack(spacing: 8) {
                    Circle()
                        .fill(ColorTheme.darkGray).frame(width: 12, height: 12)
                    Text("未達成")
                        .font(.system(size: 14)).foregroundColor(ColorTheme.black)
                }
            }
        }
    }
}

// 扇型の描画
struct PieSlice: Shape {
    let startAngle: Angle
    let endAngle: Angle

    func path(in rect: CGRect) -> Path {
        let center = CGPoint(x: rect.midX, y: rect.midY)
        let radius = min(rect.width, rect.height) / 2

        var path = Path()
        path.move(to: center)
        path.addArc(
            center: center,
            radius: radius,
            startAngle: startAngle,
            endAngle: endAngle,
            clockwise: false
        )
        return path
    }
}

#Preview {
    VStack(spacing: 20) {
        PieChartView(completed: 0)
        PieChartView(completed: 40)
        PieChartView(completed: 80)
    }
}
