//
//  MenuBarView.swift
//  SelfPomodoro
//
//  Created by し on 2025/08/22.
//

import SwiftUI

struct MenuBarView: View {
    let title: String

    var body: some View {
        Text(title)
            .font(.system(size: 28, weight: .bold))
            .foregroundColor(ColorTheme.navy)
            .frame(maxWidth: .infinity, alignment: .leading)
            .padding(.horizontal, 16)
            .padding(.top, 8)
    }
}

#Preview {
    VStack(spacing: 24) {
        MenuBarView(title: "Pomodoro")
        MenuBarView(title: "Tasks")
        MenuBarView(title: "Statistics")
        MenuBarView(title: "Profile")
    }
}
