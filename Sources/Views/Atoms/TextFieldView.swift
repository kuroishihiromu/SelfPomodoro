//
//  TextFieldView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/20.
//

import SwiftUI

struct NormalTextField: View {
    let placeholder: String
    let bgColor: Color
    let fontColor: Color
    let width: CGFloat
    let height: CGFloat
    let icon: Image?
    @Binding var text: String

    init(
        placeholder: String,
        bgColor: Color,
        fontColor: Color,
        icon: Image? = nil,
        width: CGFloat,
        height: CGFloat,
        text: Binding<String>
    ) {
        self.placeholder = placeholder
        self.bgColor = bgColor
        self.fontColor = fontColor
        self.icon = icon
        self.width = width
        self.height = height
        self._text = text
    }

    var body: some View {
        HStack(spacing: 10) {
            if let icon = icon {
                icon
                    .resizable()
                    .frame(width: 18, height: 18)
                    .foregroundColor(fontColor)
            }

            TextField(placeholder, text: $text)
                .foregroundColor(fontColor)
                .autocapitalization(.none)
                .keyboardType(.emailAddress)
        }
        .padding()
        .frame(width: width, height: height)
        .background(bgColor)
        .cornerRadius(6)
    }
}

#Preview {
    PreviewWrapper()
}

struct PreviewWrapper: View {
    @State var email = ""
    @State var password = ""

    var body: some View {
        VStack(spacing: 20) {
            NormalTextField(
                placeholder: "Your email address",
                bgColor: ColorTheme.lightGray,
                fontColor: ColorTheme.black,
                icon: Image(.mail),
                width: 350,
                height: 44,
                text: $email
            )
            NormalTextField(
                placeholder: "Enter your password",
                bgColor: ColorTheme.lightGray,
                fontColor: ColorTheme.black,
                icon: Image(.key),
                width: 350,
                height: 44,
                text: $password
            )
        }
    }
}
