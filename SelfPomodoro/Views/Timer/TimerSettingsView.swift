//
//  TimerSettingsView.swift
//  SelfPomodoro
//
//  Created by し on 2025/02/24.
//

import SwiftUI

struct TimerSettingsView: View {
    @State private var selectedTaskMinutes: Int = 25   // タスク時間：25分が初期値
    @State private var selectedRestMinutes: Int = 5      // レスト時間：5分が初期値
    @State private var selectedRounds: Int = 4           // ラウンド数：4が初期値
    @State private var navigateToCycle = false

    var body: some View {
        NavigationStack {
            VStack(spacing: 20) {
                Spacer()
                Text("タイマー設定")
                    .font(.title)
                    .padding(.top)
                Spacer()
                HStack(spacing: 20) {
                    // タスク時間用のPicker
                    VStack {
                        Text("タスク")
                            .font(.headline)
                        Picker("タスク", selection: $selectedTaskMinutes) {
                            ForEach(1...60, id: \.self) { minute in
                                Text("\(minute)分").tag(minute)
                            }
                        }
                        .pickerStyle(WheelPickerStyle())
                        .frame(width: 100, height: 100)
                    }

                    // レスト時間用のPicker
                    VStack {
                        Text("レスト")
                            .font(.headline)
                        Picker("レスト", selection: $selectedRestMinutes) {
                            ForEach(1...30, id: \.self) { minute in
                                Text("\(minute)分").tag(minute)
                            }
                        }
                        .pickerStyle(WheelPickerStyle())
                        .frame(width: 100, height: 100)
                    }

                    // ラウンド数用のPicker
                    VStack {
                        Text("ラウンド")
                            .font(.headline)
                        Picker("ラウンド", selection: $selectedRounds) {
                            ForEach(1...10, id: \.self) { round in
                                Text("\(round)回").tag(round)
                            }
                        }
                        .pickerStyle(WheelPickerStyle())
                        .frame(width: 100, height: 100)
                    }
                }
                Spacer()
                LetsTaskView(action: {
                    navigateToCycle = true
                })
                .padding(.horizontal)
                Spacer()
            }
            .padding()
            .navigationBarTitleDisplayMode(.inline)
            .navigationDestination(isPresented: $navigateToCycle) {
                CycleView(timerViewModel: TimerViewModel(
                    totalRounds: selectedRounds,
                    taskDuration: selectedTaskMinutes * 60,
                    restDuration: selectedRestMinutes * 60
                ))
            }
        }
    }
}

#Preview {
    TimerSettingsView()
}
