//
//  TextFieldView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/20.
//

import SwiftUI

struct NormalTextField: View {
    let placeholder: String
    let width: CGFloat
    let height: CGFloat
    let icon: Image?
    @Binding var text: String

    private let bgColor = ColorTheme.lightGray
    private let fontColor = ColorTheme.black

    init(
        placeholder: String,
        icon: Image? = nil,
        width: CGFloat,
        height: CGFloat,
        text: Binding<String>
    ) {
        self.placeholder = placeholder
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
    PreviewWrapperTextField()
}

struct PreviewWrapperTextField: View {
    @State var email = ""
    @State var password = ""

    var body: some View {
        VStack(spacing: 20) {
            NormalTextField(
                placeholder: "Your email address",
                icon: Image(.mail),
                width: 350,
                height: 44,
                text: $email
            )
            NormalTextField(
                placeholder: "Enter your password",
                icon: Image(.key),
                width: 350,
                height: 44,
                text: $password
            )
        }
    }
}
