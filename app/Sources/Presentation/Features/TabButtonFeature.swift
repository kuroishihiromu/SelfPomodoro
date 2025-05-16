//
//  TabButtonFeature.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/04/07.
//

import ComposableArchitecture
import Foundation

@Reducer
struct TabButtonFeature {
    @ObservableState
    struct State: Equatable {
        var selectedTabIndex: Int = 0
    }

    enum Action {
        case homeButtonTapped
        case tasksButtonTapped
        case statsButtonTapped
        case profileButtonTapped
    }

    var body: some ReducerOf<Self> {
        Reduce { state, action in
            switch action {
            case .homeButtonTapped:
                state.selectedTabIndex = 0
                return .none
            case .tasksButtonTapped:
                state.selectedTabIndex = 1
                return .none
            case .statsButtonTapped:
                state.selectedTabIndex = 2
                return .none
            case .profileButtonTapped:
                state.selectedTabIndex = 3
                return .none
            }
        }
    }
}
