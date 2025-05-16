//
//  ToDoListView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/04/19.
//

import SwiftUI
import ComposableArchitecture

struct ToDoListView: View {
    let store: StoreOf<ToDoListFeature>

    @State private var newTaskTitle = ""

    var body: some View {
        WithViewStore(store, observe: { $0 }) { viewStore in
            VStack(spacing: 16) {
                ScrollView {
                    VStack{
                        ForEachStore(
                            store.scope(state: \.items, action: \.items)
                        ) { itemStore in
                            ToDoListRow(store: itemStore, width:350, height:30)
                        }
                    }
                }
            }
            .padding()
            VStack {
                TextField("新しいタスク", text: $newTaskTitle)
                    .textFieldStyle(.roundedBorder)
                    .frame(width: 350, height: 50)

                NormalButton(
                    text: "Add Task",
                    bgColor: ColorTheme.navy,
                    fontColor: ColorTheme.white,
                    icon: Image(.add),
                    width:350,
                    height: 40,
                    action: {
                        guard !newTaskTitle.isEmpty else { return }
                        viewStore.send(.addItem(title: newTaskTitle))
                        newTaskTitle = ""
                    }
                )
            }

        }
    }
}

#Preview {
    ToDoListView(
        store: Store(initialState: ToDoListFeature.State()) {
            ToDoListFeature()
        }
    )
}
