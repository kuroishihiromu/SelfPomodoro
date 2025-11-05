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
        case addItem(detail: String)
        case addItemResponse(Result<Task, Error>)
        case items(IdentifiedActionOf<ToDoListRowFeature>)
    }
    
    @Dependency(\.taskRepository) private var taskRepository
    @Dependency(\.userIdentifier) private var userIdentifier

    var body: some ReducerOf<Self> {
        Reduce { state, action in
            switch action {
            case let .addItem(detail):
                return .run { send in
                    do {
                        let task = try await taskRepository.createTask(
                            detail: detail,
                            for: userIdentifier()
                        )
                        await send(.addItemResponse(.success(task)))
                    } catch {
                        await send(.addItemResponse(.failure(error)))
                    }
                }
                
            case let .addItemResponse(.success(task)):
                state.items.append(.init(id: task.id, detail: task.detail, isCompleted: task.isCompleted))
                return .none

            case .addItemResponse(.failure):
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
