//
//  ToDoListFeature.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/04/19.
//

import ComposableArchitecture
import SwiftUI
import Foundation

@Reducer
struct ToDoListFeature {
    struct State: Equatable {
        var items: IdentifiedArrayOf<ToDoListRowFeature.State> = []
    }

    enum Action {
        case addItem(title: String)
        case items(IdentifiedActionOf<ToDoListRowFeature>)
    }

    var body: some ReducerOf<Self> {
        Reduce { state, action in
            switch action {
            case let .addItem(title):
                state.items.append(.init(id: UUID(), title: title, isCompleted: false))
                return .none

            case .items:
                return .none
            }
        }
        .forEach(\.items, action: \.items) {
            ToDoListRowFeature()
        }
    }
}

