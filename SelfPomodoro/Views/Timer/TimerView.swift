//
//  TimerView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/01/10.
//

import SwiftUI

struct TimerView: View {
    @ObservedObject private var timerViewModel = TimerViewModel(totalRounds: 5)
    @State private var isTimeActive = false
    
    let bgColor = Color.init(red:0.90, green: 0.92, blue: 0.98)
    let buttonColor = Color.init(red: 0.38, green: 0.28, blue: 0.86)
    let lightColor = Color.init(red: 0.54, green: 0.41, blue: 0.95)
    let shadowColor = Color.init(red: 0.25, green: 0.17, blue: 0.75)
    
    var body: some View {
        VStack {
            Spacer()
            Text("\(timerViewModel.state == .task ? "タスク中" : "休憩中")")
                .font(.system(size: 20, weight: .semibold, design: .default))
                .foregroundColor(.black)
            .foregroundColor(.black)
            HStack {
                Text("\(timerViewModel.formattedTime)  /  \(timerViewModel.totalTaskDuration)")
                    .frame(width: 200, height: 200)
                    .font(.system(size: 20, weight: .semibold, design: .default))
                    .foregroundColor(.white)
                    .padding(.horizontal, 35)
                    .padding(.vertical, 20)
                    .background(
                        Circle()
                            .fill(
                                .shadow(.inner(color: lightColor, radius: 6, x: 4, y: 4))
                                .shadow(.inner(color: shadowColor, radius: 6, x: -2, y: -2))
                            )
                            .foregroundColor(buttonColor)
                            .shadow(color: buttonColor, radius: 20, y: 10)
                    )
            }
            HStack{
                Text("\(timerViewModel.round)")
                    .font(.system(size: 30, weight: .bold))
                    .foregroundColor(.black)

                Text("/ \(timerViewModel.totalRounds)")
                    .font(.system(size: 20, weight: .semibold))
                    .foregroundColor(.black)

            }
            Spacer()
        }
        StartTimerView(isActive: $isTimeActive, timerViewModel: timerViewModel)
        StopTimerView(isActive: $isTimeActive, timerViewModel: timerViewModel)
    }
}

#Preview {
    TimerView()
}
