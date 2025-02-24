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
                StartView(action: {
                    withAnimation(.easeInOut(duration: 0.3)) {
                        showSettings = true
                    }
                })
                Spacer()
            }

            if showSettings {
                // 背景が薄暗くなるように設定
                Color.black.opacity(0.4)
                    .ignoresSafeArea()
                    .onTapGesture {
                        withAnimation(.easeInOut(duration: 0.3)) {
                            showSettings = false
                        }
                    }

                // GeometryReaderを利用して画面サイズに合わせたモーダルのサイズにする
                GeometryReader { geometry in
                    TimerSettingsView()
                        .frame(width: geometry.size.width * 0.9, height: geometry.size.height * 0.5)
                        .background(Color.white)
                        .cornerRadius(20)
                        .shadow(radius: 10)
                        .transition(AnyTransition.scale.combined(with: .opacity))
                        .position(x: geometry.size.width / 2, y: geometry.size.height / 2)
                }
            }
        }
        .animation(.easeInOut(duration: 0.5), value: showSettings)
    }
}

#Preview {
    HomeView()
}
