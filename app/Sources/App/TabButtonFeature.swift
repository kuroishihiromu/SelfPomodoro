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
        var chartFeatureState: ChartFeature.State = .init()
    }

    enum Action {
        case homeButtonTapped
        case tasksButtonTapped
        case statsButtonTapped
        case profileButtonTapped
        
        case fetchTasksResponse(Result<[Task], Error>)
        case todoList(ToDoListFeature.Action)
    }
    
    
    @Dependency(\.taskRepository) private var taskRepository
    @Dependency(\.userIdentifier) private var userIdentifier

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
                return .run { send in
                    do {
                        let tasks = try await taskRepository.fetchTasks(for: userIdentifier())
                        await send(.fetchTasksResponse(.success(tasks)))
                    } catch {
                        await send(.fetchTasksResponse(.failure(error)))
                    }

                }
                
            case let .fetchTasksResponse(.success(tasks)):
                state.todoListState.items = IdentifiedArrayOf(
                    uniqueElements: tasks.map { task in
                        ToDoListRowFeature.State(
                            id: task.id,
                            detail: task.detail,
                            isCompleted: task.isCompleted
                        )
                    }
                )
                return .none

            case .fetchTasksResponse(.failure):
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
