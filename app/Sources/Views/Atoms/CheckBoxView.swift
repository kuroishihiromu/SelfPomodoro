//
//  CheckBoxView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/06/16.
//

import Foundation
import SwiftUI
struct CheckboxView: View {
    @Binding var isChecked: Bool

    var body: some View {
        Button(action: {
            isChecked.toggle()
        }) {
            HStack {
                ZStack {
                    Image(systemName: isChecked ? "checkmark.square.fill" : "square")
                        .foregroundColor(ColorTheme.navy)
                        .font(.title2)
                    if isChecked {
                        Image(systemName: "checkmark")
                            .foregroundColor(ColorTheme.white)
                            .font(.system(size: 10, weight: .bold))
                    }
                }
            }
        }
        .buttonStyle(PlainButtonStyle())
    }
}

#Preview {
    PreviewCheckbox()
}

struct PreviewCheckbox: View {
    @State private var agreed = false

    var body: some View {
        CheckboxView(isChecked: $agreed)
    }
}
