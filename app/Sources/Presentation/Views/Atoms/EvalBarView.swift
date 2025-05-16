//
//  EvalBarView.swift
//  SelfPomodoro
//
//  Created by し on 2025/03/29.
//

import SwiftUI

struct EvalBar: View {
    @Binding var value: Double
    let width: CGFloat
    let height: CGFloat

    private let leftbarColor = ColorTheme.navy
    private let rightbarColor = ColorTheme.skyBlue
    private let circleColor = ColorTheme.white
    private let circleborderColor = ColorTheme.black

    init(
        value: Binding<Double>,
        width: CGFloat,
        height: CGFloat
    ) {
        _value = value
        self.width = width
        self.height = height
    }

    var body: some View {
        ZStack(alignment: .leading) {
            // バー右側
            Capsule()
                .fill(rightbarColor)
                .frame(width: width, height: height * 0.4)

            // バー左側
            Capsule()
                .fill(leftbarColor)
                .frame(width: CGFloat(value) * width, height: height * 0.4)

            // 丸いつまみ
            Circle()
                .fill(circleColor)
                .frame(width: width, height: height)
                .overlay(
                    Circle()
                        .stroke(circleborderColor, lineWidth: 1)
                )
                .offset(x: CGFloat(value) * width - width / 2)
                .gesture(
                    DragGesture(minimumDistance: 0)
                        .onChanged { gesture in
                            let newValue = min(max(gesture.location.x / width, 0.0), 1.0)
                            value = newValue
                        }
                )
        }
        .frame(width: width, height: height)
    }
}

#Preview {
    PreviewEvalBar()
}

struct PreviewEvalBar: View {
    @State private var score: Double = 0.5

    var body: some View {
        VStack(spacing: 16) {
            EvalBar(
                value: $score,
                width: 200,
                height: 20
            )
            Text("評価: \(Int(score * 100))点")
        }
    }
}
