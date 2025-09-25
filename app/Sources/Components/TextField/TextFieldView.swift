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

struct OTPTextField: View {
    @Binding var otpCode: String
    @FocusState private var isFocused: Bool
    
    private let bgColor = ColorTheme.lightGray
    private let fontColor = ColorTheme.black
    
    var body: some View {
        VStack(spacing: 16) {
            Text("確認コードを入力してください")
                .font(.system(size: 16, weight: .medium))
                .foregroundColor(fontColor)
            
            HStack(spacing: 12) {
                ForEach(0..<6, id: \.self) { index in
                    ZStack {
                        RoundedRectangle(cornerRadius: 8)
                            .stroke(isFocused && index == otpCode.count ? Color.blue : Color.gray, lineWidth: 1)
                            .frame(width: 45, height: 50)
                            .background(bgColor)
                            .cornerRadius(8)
                        
                        Text(getDigit(at: index))
                            .font(.system(size: 20, weight: .semibold))
                            .foregroundColor(fontColor)
                    }
                }
            }
            
            TextField("", text: $otpCode)
                .keyboardType(.numberPad)
                .textContentType(.oneTimeCode)
                .focused($isFocused)
                .opacity(0.01)
                .blendMode(.screen)
                .onChange(of: otpCode) { _, newValue in
                    if newValue.count > 6 {
                        otpCode = String(newValue.prefix(6))
                    }
                    // 数字以外を除去
                    otpCode = otpCode.filter { $0.isNumber }
                }
        }
        .onTapGesture {
            isFocused = true
        }
        .onAppear {
            isFocused = true
        }
    }
    
    private func getDigit(at index: Int) -> String {
        guard index < otpCode.count else { return "" }
        return String(otpCode[otpCode.index(otpCode.startIndex, offsetBy: index)])
    }
}
