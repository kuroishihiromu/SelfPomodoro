//
//  NormalButtonView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/20.
//

import SwiftUI

struct NormalButton: View {
    let text: String
    let bgColor: Color
    let fontColor: Color
    let width: CGFloat
    let height: CGFloat
    let icon: Image?
    let action: () -> Void

    init(
        text: String,
        bgColor: Color,
        fontColor: Color,
        icon: Image? = nil,
        width: CGFloat,
        height: CGFloat,
        action: @escaping () -> Void
    ) {
        self.text = text
        self.bgColor = bgColor
        self.fontColor = fontColor
        self.icon = icon
        self.width = width
        self.height = height
        self.action = action
    }

    var body: some View {
        Button(action: action) {
            HStack(spacing: 15) {
                if let icon = icon {
                    icon
                        .resizable()
                        .frame(width: 18, height: 18)
                }

                Text(text)
                    .foregroundColor(fontColor)
                    .font(.headline)
            }
            .padding()
            .frame(width: width, height: height)
            .background(bgColor)
            .cornerRadius(10)
        }
    }
}
