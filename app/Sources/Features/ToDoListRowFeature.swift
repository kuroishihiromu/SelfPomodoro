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
        case toggleCompletedResponse(Result<TaskResult, taskAPIError>)
    }

    @Dependency(\.taskAPIClient) var apiClient
    
    func reduce(into state: inout State, action: Action) -> Effect<Action> {
        switch action {
        case .toggleCompleted:
           print("🟢 toggleCompleted called for id: \(state.id), current: \(state.isCompleted)")
           return .run { [id = state.id] send in
               do {
                   let result = try await apiClient.toggleCompletion(id)
                   print("✅ toggleCompletion succeeded: \(result)")
                   await send(.toggleCompletedResponse(.success(result)))
               } catch let apiError as taskAPIError {
                   print("❌ toggleCompletion failed with taskAPIError: \(apiError)")
                   await send(.toggleCompletedResponse(.failure(apiError)))
               } catch {
                   print("❌ toggleCompletion failed with unknown error: \(error)")
                   await send(.toggleCompletedResponse(.failure(.unknown)))
               }
           }
        case let .toggleCompletedResponse(.success(result)):
            print("🟡 toggleCompletedResponse (success): \(result)")
            state.isCompleted = result.isCompleted
            return .none

        case let .toggleCompletedResponse(.failure(err)):
            print("🔴 toggleCompletedResponse (failure): \(err)")
            return .none
        }
    }
}
