//
//  SettingTimerField.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/20.
//

import SwiftUI

struct SettingTimerField: View {
    let placeholder: String
    let width: CGFloat
    let height: CGFloat
    @Binding var text: String

    private let bgColor = ColorTheme.white
    private let borderColor = ColorTheme.Gray
    private let fontColor = ColorTheme.black

    init(
        placeholder: String,
        width: CGFloat,
        height: CGFloat,
        text: Binding<String>
    ) {
        self.placeholder = placeholder
        self.width = width
        self.height = height
        self._text = text
    }

    var body: some View {
        HStack{
            TextField(placeholder, text: $text)
                .foregroundColor(fontColor)
                .keyboardType(.numberPad)
        }
        .padding()
        .frame(width: width, height: height)
        .background(bgColor)
        .overlay(
            RoundedRectangle(cornerRadius: 6)
                .stroke(borderColor)
        )
    }
}

#Preview {
    PreviewWrapperSettingTimerField()
}

struct PreviewWrapperSettingTimerField: View {
    @State var time: String = ""

    var body: some View {
        VStack {
            SettingTimerField(
                placeholder: "25",
                width: 160,
                height: 44,
                text: $time
            )
        }
    }
}
