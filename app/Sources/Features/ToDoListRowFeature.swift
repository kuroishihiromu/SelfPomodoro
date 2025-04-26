//
//  ToDoListRowFeature.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/04/19.
//
import ComposableArchitecture
import Foundation

@Reducer
struct ToDoListRowFeature {
    struct State: Equatable, Identifiable {
        let id: UUID
        var title: String
        var isCompleted: Bool
    }

    enum Action {
        case toggleCompleted
    }

    func reduce(into state: inout State, action: Action) -> Effect<Action> {
        switch action {
        case .toggleCompleted:
            state.isCompleted.toggle()
            return .none
        }
    }
}
