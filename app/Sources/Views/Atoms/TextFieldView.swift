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
        _text = text
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


struct PasswordTextField: View {
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
        _text = text
    }

    var body: some View {
        HStack(spacing: 10) {
            if let icon = icon {
                icon
                    .resizable()
                    .frame(width: 18, height: 18)
                    .foregroundColor(fontColor)
            }

            SecureField(placeholder, text: $text)
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
