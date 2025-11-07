//
//  StatisticsScreenView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/20.
//

import SwiftUI
import ComposableArchitecture

@Reducer
struct StatisticsFeature {
    enum DisplayMode: Equatable {
        case chart
        case heatMap
    }

    @ObservableState
    struct State: Equatable {
        var chart = ChartFeature.State()
        var heatMap = HeatMapFeature.State()
        var displayMode: DisplayMode = .chart
        var totalCompletedRounds: Int = 0
        var totalWorkMinutes: Int = 0
    }
    
    enum Action {
        case chart(ChartFeature.Action)
        case heatMap(HeatMapFeature.Action)
        case changeMode(DisplayMode)
        case fetchTotalRounds
        case fetchTotalRoundsResponse(Result<Int, Error>)
        case fetchTotalWork
        case fetchTotalWorkResponse(Result<Int, Error>)
    }
    
    @Dependency(\.roundRecordRepository) var roundRecordRepository
    @Dependency(\.userIdentifier) var userIdentifier

    var body: some ReducerOf<Self> {
        Scope(state: \.chart, action: \.chart) {
            ChartFeature()
        }
        Scope(state: \.heatMap, action: \.heatMap) {
            HeatMapFeature()
        }
        Reduce { state, action in
            switch action {
            case .changeMode(let mode):
                state.displayMode = mode
                return .none
            case .chart, .heatMap:
                return .none
            case .fetchTotalRounds:
                return .run { send in
                    let identifier = userIdentifier()
                    let count = try await roundRecordRepository.countCompleted(for: identifier)
                    await send(.fetchTotalRoundsResponse(.success(count)))
                } catch: { error, send in
                    await send(.fetchTotalRoundsResponse(.failure(error)))
                }
            case let .fetchTotalRoundsResponse(.success(count)):
                state.totalCompletedRounds = count
                return .none
            case .fetchTotalRoundsResponse(.failure):
                return .none
            case .fetchTotalWork:
                return .run { send in
                    let identifier = userIdentifier()
                    let minutes = try await roundRecordRepository.totalWorkMinutes(for: identifier)
                    await send(.fetchTotalWorkResponse(.success(Int(minutes.rounded()))))
                } catch: { error, send in
                    await send(.fetchTotalWorkResponse(.failure(error)))
                }
            case let .fetchTotalWorkResponse(.success(total)):
                state.totalWorkMinutes = total
                return .none
            case .fetchTotalWorkResponse(.failure):
                return .none
        }
    }
}
}

struct StatisticsScreenView: View {
    let store: StoreOf<StatisticsFeature>
    @State private var animatedRounds: Int = 0
    @State private var animatedWorkMinutes: Int = 0
    @State private var roundsAnimationTask: Task<Void, Never>? = nil

    var body: some View {
        WithViewStore(store, observe: \.self) { viewStore in
            VStack {
                // 切り替えボタン
                SegmentedTabView(
                    selection: viewStore.binding(get: { $0.displayMode }, send: { .changeMode($0) })
                )
                .padding(.horizontal)
                .padding(.top, 12)

                // チャートかヒートマップ
                VStack {
                    if viewStore.displayMode == .chart {
                        ChartView(store: store.scope(state: \.chart, action: \.chart))
                            .transition(.opacity)
                    } else {
                        HeatMapView(store: store.scope(state: \.heatMap, action: \.heatMap))
                            .transition(.opacity)
                    }
                    // 累計ラウンド数
                    VStack(spacing: 8) {
                        Text("達成したラウンド数")
                            .font(.headline)
                            .foregroundColor(ColorTheme.black)
                        Text(String(animatedRounds))
                            .font(.system(size: 30))
                            .monospacedDigit()
                    }
                    .padding(.top, 50)

                    // 総作業時間（分）
                    VStack(spacing: 8) {
                        Text("総作業時間(分)")
                            .font(.headline)
                            .foregroundColor(ColorTheme.black)
                        Text(String(animatedWorkMinutes))
                            .font(.system(size: 30))
                            .monospacedDigit()
                    }
                    .padding(.top, 16)
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .center)
                .animation(.easeInOut, value: viewStore.displayMode)
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)
            .safeAreaInset(edge: .top, spacing: 0) {
                MenuBarView(title: "Statistics")
            }
            .onAppear {
                viewStore.send(.fetchTotalRounds)
                viewStore.send(.fetchTotalWork)
                // 値が変わらない場合でも毎回アピアランスで再度アニメ開始
                startSynchronizedAnimation(
                    roundsTarget: viewStore.totalCompletedRounds,
                    workTarget: viewStore.totalWorkMinutes
                )
            }
            .onChange(of: viewStore.totalCompletedRounds) { _, newRounds in
                startSynchronizedAnimation(roundsTarget: newRounds, workTarget: viewStore.totalWorkMinutes)
            }
            .onChange(of: viewStore.totalWorkMinutes) { _, newWork in
                startSynchronizedAnimation(roundsTarget: viewStore.totalCompletedRounds, workTarget: newWork)
            }
            .onDisappear {
                roundsAnimationTask?.cancel()
                roundsAnimationTask = nil
            }
        }
    }
}

// MARK: - Animation Helper
extension StatisticsScreenView {
    private func startSynchronizedAnimation(roundsTarget: Int, workTarget: Int) {
        let rounds = max(0, roundsTarget)
        let work = max(0, workTarget)

        roundsAnimationTask?.cancel()
        if rounds == 0 {
            animatedRounds = 0
            animatedWorkMinutes = work
            return
        }
        animatedRounds = 0
        animatedWorkMinutes = 0

        roundsAnimationTask = Task {
            // 12ms間隔で、ラウンドの増分に同期して作業時間も線形に進める
            for i in 0...rounds {
                try? await Task.sleep(nanoseconds: 12_000_000)
                let syncedWork = Int((Double(work) * Double(i) / Double(rounds)).rounded(.down))
                await MainActor.run {
                    animatedRounds = i
                    animatedWorkMinutes = syncedWork
                }
            }
            // 最終値を正に揃える（丸め誤差ケア）
            await MainActor.run {
                animatedRounds = rounds
                animatedWorkMinutes = work
            }
        }
    }
}
