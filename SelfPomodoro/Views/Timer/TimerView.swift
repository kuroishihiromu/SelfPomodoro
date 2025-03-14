//
//  TimerView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/01/10.
//

import SwiftUI

struct TimerView: View {
    @ObservedObject var timerViewModel = TimerViewModel(totalRounds: 5, taskDuration: 1500, restDuration: 300)

    let bgColor = Color.init(red:0.90, green: 0.92, blue: 0.98)

    private var taskColor: (circle: Color, light: Color, shadow: Color) {
        return timerViewModel.state == .task
            ? (Color(red: 0.38, green: 0.28, blue: 0.86),   // タスク時のcircle色
               Color(red: 0.54, green: 0.41, blue: 0.95),   // タスク時のlight色
               Color(red: 0.25, green: 0.17, blue: 0.75))   // タスク時のshadow色
            : (Color(red: 0.38, green: 0.85, blue: 0.28),   // 休憩時のcircle色
               Color(red: 0.54, green: 0.95, blue: 0.41),   // 休憩時のlight色
               Color(red: 0.25, green: 0.75, blue: 0.17))   // 休憩時のshadow色
    }

    private var circleColor: Color { taskColor.circle }
    private var lightColor: Color { taskColor.light }
    private var shadowColor: Color { taskColor.shadow }

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
                    .font(.system(size: 20, weight: .semibold))
                    .foregroundColor(.white)
                    .padding(.horizontal, 35)
                    .padding(.vertical, 20)
                    .background(
                        Circle()
                            .fill(
                                .shadow(.inner(color: lightColor, radius: 6, x: 4, y: 4))
                                .shadow(.inner(color: shadowColor, radius: 6, x: -2, y: -2))
                            )
                            .foregroundColor(circleColor)
                            .shadow(color: circleColor, radius: 20, y: 10)
                    )
            }
            .padding(.bottom, 50)
            HStack{
                Text("\(timerViewModel.round)")
                    .font(.system(size: 40, weight: .bold))
                    .foregroundColor(.black)

                Text("/ \(timerViewModel.totalRounds)")
                    .font(.system(size: 30, weight: .semibold))
                    .foregroundColor(.black)

            }
            Spacer()
            if timerViewModel.isActive {
                StopTimerView(action: {
                    timerViewModel.stopTimer()
                })
            } else {
                StartTimerView(action: {
                    timerViewModel.startTimer()
                })
            }
        }
    }
}

#Preview {
    TimerView()
}
