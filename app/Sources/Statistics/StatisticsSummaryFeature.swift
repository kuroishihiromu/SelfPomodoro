//
//  StatisticsSummaryFeature.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/08.
//

import ComposableArchitecture

@Reducer
struct StatisticsSummaryFeature {

    @ObservableState
    struct State: Equatable {
        var totalCompletedRounds: Int = 0
        var totalWorkMinutes: Int = 0
    }

    enum Action: Equatable {
        case onAppear
        case totalRoundsResponse(Int)
        case totalWorkResponse(Int)
    }

    @Dependency(\.roundRecordRepository) var roundRecordRepository
    @Dependency(\.userIdentifier) var userIdentifier

    var body: some ReducerOf<Self> {
        Reduce { state, action in
            switch action {
            case .onAppear:
                let identifier = userIdentifier()
                print("📊 StatisticsSummaryFeature.onAppear for user=\(identifier)")
                return .merge(
                    .run { send in
                        let count = try await roundRecordRepository.countCompleted(for: identifier)
                        print("📊 Summary totalRounds fetched=\(count)")
                        await send(.totalRoundsResponse(count))
                    } catch: { error, send in
                        print("⚠️ StatisticsSummaryFeature: round count fetch failed \(error)")
                    },
                    .run { send in
                        let minutes = try await roundRecordRepository.totalWorkMinutes(for: identifier)
                        let rounded = Int(minutes.rounded())
                        print("📊 Summary totalWorkMinutes fetched=\(rounded)")
                        await send(.totalWorkResponse(rounded))
                    } catch: { error, send in
                        print("⚠️ StatisticsSummaryFeature: work minutes fetch failed \(error)")
                    }
                )

            case let .totalRoundsResponse(count):
                print("📊 Summary totalRoundsResponse apply count=\(count)")
                state.totalCompletedRounds = count
                return .none

            case let .totalWorkResponse(total):
                print("📊 Summary totalWorkResponse apply minutes=\(total)")
                state.totalWorkMinutes = total
                return .none
            }
        }
    }
}
