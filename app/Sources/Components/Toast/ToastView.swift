//
//  ToastView.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/08.
//

import SwiftUI

struct ToastState: Equatable {
    var message: String
    var isVisible: Bool = false
    var tint: Color = ColorTheme.navy
}

struct ToastView: View {
    let state: ToastState

    var body: some View {
        if state.isVisible {
            Text(state.message)
                .font(.caption)
                .foregroundColor(.white)
                .padding(.horizontal, 16)
                .padding(.vertical, 10)
                .background(RoundedRectangle(cornerRadius: 12).fill(state.tint))
                .transition(.move(edge: .top).combined(with: .opacity))
                .padding(.top, 16)
        }
    }
}

