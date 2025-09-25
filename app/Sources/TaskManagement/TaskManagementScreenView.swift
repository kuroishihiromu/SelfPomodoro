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
        VStack {
            ToDoListView(store: store)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .safeAreaInset(edge: .top, spacing: 0) {
            MenuBarView(title: "Tasks")
        }
    }
}

#Preview {
    TaskManagementScreenView(
        store: Store(initialState: ToDoListFeature.State()) {
            ToDoListFeature()
        }
    )
}
