//
//  HomeView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/01/10.
//

import SwiftUI

struct HomeView: View {
    @State private var showSettings = false  // モーダル表示の有無を管理する状態変数

    var body: some View {
        ZStack {
            VStack {
                Spacer()
                Button(action: {
                    showSettings = true
                }) {
                    Text("Start")
                        .font(.largeTitle)
                        .foregroundColor(.white)
                        .padding()
                        .background(Color(.blue))
                        .cornerRadius(10)
                }
                Spacer()
            }
        }
        .sheet(isPresented: $showSettings) {  // スタートボタンが押されたときにモーダルを表示
//            TimerSettingsView()
        }
    }
}

#Preview {
    HomeView()
}
