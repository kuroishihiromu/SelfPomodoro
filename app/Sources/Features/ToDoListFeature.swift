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
        case addItemResponse(Result<TaskResult, taskAPIError>)
        case items(IdentifiedActionOf<ToDoListRowFeature>)
    }
    
    @Dependency(\.taskAPIClient) var apiClient

    var body: some ReducerOf<Self> {
        Reduce { state, action in
            switch action {
            case let .addItem(title):
                print("🟢 addItem called with title: \(title)")
                state.items.append(.init(id: UUID(), title: title, isCompleted: false))
                return .run { send in
                    do {
                        let result = try await apiClient.addTask(title)
                        print("✅ addTask succeeded: \(result)")
                        await send(.addItemResponse(.success(result)))
                    } catch let error as taskAPIError {
                        print("❌ addTask failed with taskAPIError: \(error)")
                        await send(.addItemResponse(.failure(error)))
                    } catch {
                        print("❌ addTask failed with unknown error: \(error)")
                        await send(.addItemResponse(.failure(.unknown)))
                    }
                }
                
            case let .addItemResponse(.success(task)):
                print("🟡 addItemResponse (success): \(task)")
                state.items.append(.init(id: task.id, title: task.detail, isCompleted: task.isCompleted))
                return .none

            case let .addItemResponse(.failure(toggleCompleteResponseErr)):
                print("🔴 addItemResponse (failure): \(toggleCompleteResponseErr)")
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

