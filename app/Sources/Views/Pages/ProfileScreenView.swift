//
//  ProfileScreenView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/02/27.
//

import SwiftUI

struct ProfileScreenView: View {
    var body: some View {
        VStack {
            Text("coming soon...")
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .safeAreaInset(edge: .top, spacing: 0) {
            MenuBarView(title: "Profile")
        }
    }
}

#Preview {
    ProfileScreenView()
}
