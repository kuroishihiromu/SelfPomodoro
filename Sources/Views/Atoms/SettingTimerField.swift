//
//  SettingTimerField.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/20.
//

import SwiftUI

struct SettingTimerButton: View {
    let width: CGFloat
    let height: CGFloat
    @Binding var duration: TimeInterval

    @State private var showPicker = false

    private let bgColor = ColorTheme.white
    private let borderColor = ColorTheme.Gray
    private let fontColor = ColorTheme.black

    var body: some View {
        Button(action: {
            showPicker.toggle()
        }) {
            Text(formatTime(duration))
                .foregroundColor(fontColor)
                .frame(width: width, height: height)
                .background(bgColor)
                .overlay(
                    RoundedRectangle(cornerRadius: 6)
                        .stroke(borderColor)
                )
        }
        .sheet(isPresented: $showPicker) {
            TimePickerSheet(duration: $duration)
        }
    }

    private func formatTime(_ interval: TimeInterval) -> String {
        let minutes = Int(interval) / 60
        return String(format: "%02d:%02d", minutes, Int(interval) % 60)
    }
}

struct TimePickerSheet: View {
    @Binding var duration: TimeInterval
    @Environment(\.dismiss) var dismiss

    @State private var minutes: Int
    @State private var seconds: Int

    init(duration: Binding<TimeInterval>) {
        self._duration = duration
        let totalSeconds = Int(duration.wrappedValue)
        _minutes = State(initialValue: totalSeconds / 60)
        _seconds = State(initialValue: totalSeconds % 60)
    }

    var body: some View {
        VStack(spacing: 20) {
            Text("時間を選択")
                .font(.headline)

            HStack {
                Picker("分", selection: $minutes) {
                    ForEach(0..<60) { Text("\($0) 分") }
                }
                .pickerStyle(.wheel)
                .frame(maxWidth: .infinity)

                Picker("秒", selection: $seconds) {
                    ForEach(0..<60) { Text("\($0) 秒") }
                }
                .pickerStyle(.wheel)
                .frame(maxWidth: .infinity)
            }
            .frame(height: 150)

            Button("完了") {
                duration = TimeInterval(minutes * 60 + seconds)
                dismiss()
            }
            .padding()
        }
        .padding()
    }
}

#Preview {
    PreviewWrapperSettingTimerButton()
}

struct PreviewWrapperSettingTimerButton: View {
    @State var time: TimeInterval = 25 * 60  // 初期値：25分

    var body: some View {
        VStack {
            SettingTimerButton(
                width: 160,
                height: 44,
                duration: $time
            )
        }
    }
}
