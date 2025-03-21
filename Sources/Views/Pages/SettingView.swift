//
//  HomeView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/01/10.
//

import SwiftUI

struct HomeView: View {
    @State private var showSettings = false             // タイマー設定のモーダルの表示を管理
    @State private var isPresentingTimerView = false    // CycleViewへの遷移を管理
    @State private var timerViewModel: TimerViewModel?  // タイマーの設定情報を保持

    var body: some View {
        NavigationStack {
            ZStack {
                VStack {
                    Spacer()
                    // Startボタン
//                    StartView(action: {
//                        withAnimation(.easeInOut(duration: 0.3)) {
//                            showSettings = true
//                        }
//                    })
                    Spacer()
                }
                .navigationDestination(isPresented: $isPresentingTimerView) {
                    if let viewModel = timerViewModel {
                        CycleView(timerViewModel: viewModel)
                    }
                }

                // モーダル表示
                if showSettings {
                    // グレーの背景
                    Color.black.opacity(0.4)
                        .ignoresSafeArea()
                        .onTapGesture {
                            withAnimation(.easeInOut(duration: 0.3)) {
                                showSettings = false
                            }
                        }

                    // モーダル
                    GeometryReader { geometry in  // GeometryReaderを利用して画面サイズに合わせたサイズにする
                        TimerSettingsView(
                            onSettingsConfirmed: { rounds, taskTime, restTime in
                                showSettings = false
                                DispatchQueue.main.asyncAfter(deadline: .now() + 0.3) {
                                    timerViewModel = TimerViewModel(
                                        totalRounds: rounds,
                                        taskDuration: taskTime * 60,
                                        restDuration: restTime * 60
                                    )
                                    isPresentingTimerView = true
                                }
                            }
                        )
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
}

#Preview {
    HomeView()
}
