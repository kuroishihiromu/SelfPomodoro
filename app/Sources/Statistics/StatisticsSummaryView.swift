//
//  StatisticsSummaryView.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/08.
//

import SwiftUI
import ComposableArchitecture

struct StatisticsSummaryView: View {
    let store: StoreOf<StatisticsSummaryFeature>

    @State private var animatedRounds: Int = 0
    @State private var animatedWorkMinutes: Int = 0
    @State private var roundsAnimationTask: Task<Void, Never>? = nil

    var body: some View {
        WithViewStore(store, observe: { $0 }) { viewStore in
            ScrollView {
                VStack(spacing: 20) {
                    statCard(
                        title: "達成したラウンド数",
                        value: animatedRounds,
                        description: "中断されなかった完了ラウンドのみを集計",
                        systemImage: "target"
                    )

                    statCard(
                        title: "総作業時間(分)",
                        value: animatedWorkMinutes,
                        description: "ラウンドの作業フェーズ合計時間",
                        systemImage: "clock.arrow.circlepath"
                    )
                }
                .padding(.horizontal, 20)
                .padding(.vertical, 24)
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)
            .background(ColorTheme.Gray.opacity(0.1))
            .onAppear {
                viewStore.send(.onAppear)
                startSynchronizedAnimation(
                    roundsTarget: viewStore.totalCompletedRounds,
                    workTarget: viewStore.totalWorkMinutes
                )
            }
            .onChange(of: viewStore.totalCompletedRounds) { _, newValue in
                startSynchronizedAnimation(roundsTarget: newValue, workTarget: viewStore.totalWorkMinutes)
            }
            .onChange(of: viewStore.totalWorkMinutes) { _, newValue in
                startSynchronizedAnimation(roundsTarget: viewStore.totalCompletedRounds, workTarget: newValue)
            }
            .onDisappear {
                roundsAnimationTask?.cancel()
                roundsAnimationTask = nil
            }
        }
    }
}

// MARK: - Animation Helper
private extension StatisticsSummaryView {
    @ViewBuilder
    func statCard(title: String, value: Int, description: String, systemImage: String) -> some View {
        VStack(alignment: .leading, spacing: 12) {
            Label(title, systemImage: systemImage)
                .font(.headline)
                .foregroundColor(ColorTheme.black)

            Text("\(value)")
                .font(.system(size: 40, weight: .bold))
                .monospacedDigit()
                .foregroundColor(ColorTheme.navy)

            Text(description)
                .font(.caption)
                .foregroundColor(ColorTheme.black.opacity(0.6))
        }
        .padding(20)
        .frame(maxWidth: .infinity, alignment: .leading)
        .background(
            RoundedRectangle(cornerRadius: 18, style: .continuous)
                .fill(ColorTheme.white)
                .shadow(color: Color.black.opacity(0.08), radius: 16, x: 0, y: 8)
        )
    }

    func startSynchronizedAnimation(roundsTarget: Int, workTarget: Int) {
        let rounds = max(0, roundsTarget)
        let work = max(0, workTarget)

        roundsAnimationTask?.cancel()
        roundsAnimationTask = nil

        guard rounds > 0 else {
            animatedRounds = 0
            animatedWorkMinutes = work
            return
        }

        animatedRounds = 0
        animatedWorkMinutes = 0

        roundsAnimationTask = Task {
            for i in 0...rounds {
                try? await Task.sleep(nanoseconds: 12_000_000)
                let syncedWork = Int((Double(work) * Double(i) / Double(rounds)).rounded(.down))
                await MainActor.run {
                    animatedRounds = i
                    animatedWorkMinutes = syncedWork
                }
            }
            await MainActor.run {
                animatedRounds = rounds
                animatedWorkMinutes = work
            }
        }
    }
}
