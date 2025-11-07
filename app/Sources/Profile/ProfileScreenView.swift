//
//  ProfileScreenView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/02/27.
//

import SwiftUI

struct ProfileScreenView: View {
    var body: some View {
        VStack(spacing: 30) {
            Text("coming soon...")
            Spacer()
        }
        .padding()
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .safeAreaInset(edge: .top, spacing: 0) {
            MenuBarView(title: "Profile")
        }
    }
}

#Preview {
    ProfileScreenView()
}
