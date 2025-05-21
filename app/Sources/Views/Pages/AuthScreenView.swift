//
//  AuthScreenView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/02/27.
//

import SwiftUI

struct AuthScreenView: View {
    @Environment(\.isAuthenticated) private var isAuthenticated

    var body: some View {
        Text("Auth View")
    }
}
