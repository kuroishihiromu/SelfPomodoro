//
//  ProfileScreenView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/02/27.
//

import SwiftUI
import ComposableArchitecture

struct ProfileScreenView: View {
    let authStore: StoreOf<AuthFeature>
    
    var body: some View {
        VStack(spacing: 30) {
            Text("coming soon...")
            
            Spacer()
            
            NormalButton(
                text: L10n.Profile.logout,
                bgColor: .red,
                fontColor: .white,
                width: 250,
                height: 50,
                action: {
                    authStore.send(.tappedSignOut)
                }
            )
        }
        .padding()
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .safeAreaInset(edge: .top, spacing: 0) {
            MenuBarView(title: "Profile")
        }
    }
}

#Preview {
    ProfileScreenView(
        authStore: Store(
            initialState: AuthFeature.State(),
            reducer: { AuthFeature() }
        )
    )
}
