//
//  TabButtonView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/04/09.
//

import SwiftUI

struct TabButton: View {
    let icon: String
    let isSelected: Bool
    let action: () -> Void

    @State private var animateTrigger: Bool = false
    @State private var wasSelected: Bool = false

    var body: some View {
        Image(systemName: icon)
            .font(.title2)
            .foregroundColor(isSelected ? ColorTheme.navy : ColorTheme.black)
            .symbolEffect(.bounce, options: .nonRepeating, value: animateTrigger)
            .frame(maxWidth: .infinity)
            .onChange(of: isSelected) {
                if isSelected && !wasSelected {
                    triggerBounce()
                }
                wasSelected = isSelected
            }
            .onTapGesture {
                action()
            }
    }

    private func triggerBounce() {
        animateTrigger = false
        DispatchQueue.main.asyncAfter(deadline: .now() + 0.01) {
            animateTrigger = true
        }
    }
}
