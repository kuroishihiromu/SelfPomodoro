//
//  HomeView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/01/10.
//

import SwiftUI

struct HomeView: View {
    @State private var showSettings = false

    var body: some View {
        ZStack {
            VStack {
                Spacer()
                Button(action: {
                    withAnimation(.easeInOut(duration: 0.3)) {
                        showSettings = true
                    }
                }) {
                    Text("Start")
                        .font(.largeTitle)
                        .foregroundColor(.white)
                        .padding()
                        .background(Color.blue)
                        .cornerRadius(10)
                }
                Spacer()
            }

            if showSettings {
                // 背景を半透明の黒でオーバーレイ
                Color.black.opacity(0.4)
                    .ignoresSafeArea()
                    .onTapGesture {
                        withAnimation(.easeInOut(duration: 0.3)) {
                            showSettings = false
                        }
                    }

                // 中央に配置するカスタムモーダル
                TimerSettingsView()
                    .frame(width: 300, height: 400)
                    .background(Color.white)
                    .cornerRadius(20)
                    .shadow(radius: 10)
                // スケールとフェードの組み合わせでアニメーション
                    .transition(AnyTransition.scale.combined(with: .opacity))
            }
        }
        // showSettingsの変化に合わせてアニメーションを適用
        .animation(.easeInOut(duration: 0.3), value: showSettings)
    }
}

#Preview {
    HomeView()
}
