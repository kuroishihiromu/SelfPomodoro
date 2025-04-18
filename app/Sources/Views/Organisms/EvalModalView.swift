//
//  EvalModalView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/20.
//

import SwiftUI

struct EvalModalView: View {
    @State private var score: Double = 0.5
    var body: some View {
        ZStack{
            VStack{
                Text("タスクの評価")
                    .font(.system(size: 25, weight: .bold))
                    .padding(.bottom, 30)
                Text("Round 1 / 5")
                    .font(.system(size: 22, weight: .bold))
                    .padding(.bottom, 5)
                VStack{
                    Text("\(Int(score * 100)) / 100")
                        .font(.system(size: 30, weight: .bold))
                    EvalBar(
                        value: $score,
                        width: 200,
                        height: 20
                    )
                }
                .padding(.bottom, 30)
                NormalButton(
                    text: "Start Break",
                    bgColor: ColorTheme.navy,
                    fontColor: ColorTheme.white,
                    width: 180,
                    height: 48,
                    action: { print("break started!") }
                )
            }
        }
        .padding(EdgeInsets(top:25, leading: 70, bottom: 35, trailing: 70))

        .overlay(
            RoundedRectangle(cornerRadius: 3)
                .stroke(ColorTheme.navy, lineWidth: 4)
        )
        .frame(width: 650, height: 650)
    }
}

#Preview {
    EvalModalView()
}
