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
        var detail: String
        var isCompleted: Bool
    }

    enum Action {
        case toggleCompleted
        case toggleCompletedResponse(Result<TodoTask, Error>)
    }

    @Dependency(\.taskRepository) private var taskRepository
    
    func reduce(into state: inout State, action: Action) -> Effect<Action> {
        switch action {
        case .toggleCompleted:
           return .run { [id = state.id] send in
               do {
                   let task = try await taskRepository.toggleTaskCompletion(for: id)
                   await send(.toggleCompletedResponse(.success(task)))
               } catch {
                   await send(.toggleCompletedResponse(.failure(error)))
               }
           }
        case let .toggleCompletedResponse(.success(task)):
            state.detail = task.detail
            state.isCompleted = task.isCompleted
            return .none

        case .toggleCompletedResponse(.failure):
            return .none
        }
    }
}
