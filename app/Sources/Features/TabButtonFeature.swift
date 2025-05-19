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
        var todoListState: ToDoListFeature.State = .init()
    }

    enum Action {
        case homeButtonTapped
        case tasksButtonTapped
        case statsButtonTapped
        case profileButtonTapped
        
        case fetchTasksResponse(Result<[TaskResult], taskAPIError>)
        case todoList(ToDoListFeature.Action)
    }
    
    
    @Dependency(\.taskAPIClient) var apiClient

    var body: some ReducerOf<Self> {
        Scope(state: \.todoListState, action: \..todoList) {
            ToDoListFeature()
        }
        
        Reduce { state, action in
            switch action {
            case .homeButtonTapped:
                state.selectedTabIndex = 0
                return .none
            case .tasksButtonTapped:
                state.selectedTabIndex = 1
                print("📲 tasksButtonTapped → fetching tasks from API")
                return .run { send in
                    do {
                        let tasks = try await apiClient.fetchTasks()
                        print("✅ fetchTasks succeeded with \(tasks.count) items")
                        await send(.fetchTasksResponse(.success(tasks)))
                    } catch let error as taskAPIError {
                        print("❌ fetchTasks failed: \(error)")
                        await send(.fetchTasksResponse(.failure(error)))
                    } catch {
                        await send(.fetchTasksResponse(.failure(.unknown)))
                    }
                }
                
            case let .fetchTasksResponse(.success(tasks)):
                print("🟡 fetchTasksResponse (success)")
                state.todoListState.items = IdentifiedArrayOf(
                    uniqueElements: tasks.map { task in
                        ToDoListRowFeature.State(
                            id: task.id,
                            title: task.detail,
                            isCompleted: task.isCompleted
                        )
                    }
                )
                return .none

            case .fetchTasksResponse(.failure(let error)):
                print("🔴 fetchTasksResponse (failure): \(error)")
                // エラー状態に応じた UI 対応も可能
                return .none
                
            case .todoList:
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
