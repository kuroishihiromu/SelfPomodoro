//
//  TaskManagementScreenView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/20.
//

import SwiftUI
import ComposableArchitecture

struct TaskManagementScreenView: View {
    let store: StoreOf<ToDoListFeature>

    var body: some View {
        ToDoListView(store: store)
    }
}

#Preview {
    TaskManagementScreenView(
        store: Store(initialState: ToDoListFeature.State()) {
            ToDoListFeature()
        }
    )
}
