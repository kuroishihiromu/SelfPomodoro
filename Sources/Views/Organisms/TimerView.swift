//
//  TimerView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/01/10.
//

import SwiftUI

struct TimerView: View {
    @ObservedObject var timerViewModel = TimerViewModel(totalRounds: 5, taskDuration: 1500, restDuration: 300)

    let bgColor = ColorTheme.white

    private var timerColor: Color {
        return timerViewModel.state == .task ? ColorTheme.navy : ColorTheme.lightSkyBlue
    }

    private var circleColor: Color { ColorTheme.darkGray }

    var body: some View {
        VStack {
            Spacer()
            Text("\(timerViewModel.state == .task ? "タスク中" : "休憩中")")
                .font(.system(size: 20, weight: .semibold, design: .default))
                .foregroundColor(.black)

            HStack {
                Text("\(timerViewModel.formattedTime)  /  \(timerViewModel.totalTaskDuration)")
                    .frame(width: 200, height: 200)
                    .font(.system(size: 20, weight: .semibold))
                    .padding(.horizontal, 35)
                    .padding(.vertical, 20)
                    .foregroundColor(circleColor)
            }
            .padding(.bottom, 50)
            
            HStack {
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
