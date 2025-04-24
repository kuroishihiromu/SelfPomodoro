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
                let center = CGPoint(x: geometry.size.width / 2, y: geometry.size.height / 2)
                let radius = min(geometry.size.width, geometry.size.height) / 2
                let clampedCompleted = max(0, min(completed, 100)) // 0~100の範囲に収まるように処理
                let completedAngle = 360 * clampedCompleted / 100

                ZStack {
                    // 達成済み
                    PieSliceView(
                        center: center,
                        radius: radius,
                        startAngle: .degrees(-90),
                        endAngle: .degrees(completedAngle - 90),
                        color: ColorTheme.navy
                    )
                    // 未達成
                    PieSliceView(
                        center: center,
                        radius: radius,
                        startAngle: .degrees(completedAngle - 90),
                        endAngle: .degrees(270),
                        color: ColorTheme.darkGray
                    )
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
struct PieSliceView: View {
    let center: CGPoint
    let radius: CGFloat
    let startAngle: Angle
    let endAngle: Angle
    let color: Color

    var body: some View {
        Path { path in
            path.move(to: center)
            path.addArc(
                center: center,
                radius: radius,
                startAngle: startAngle,
                endAngle: endAngle,
                clockwise: false
            )
        }
        .fill(color)
    }
}

#Preview {
    VStack(spacing: 20) {
        PieChartView(completed: 0)
        PieChartView(completed: 40)
        PieChartView(completed: 80)
    }
}
