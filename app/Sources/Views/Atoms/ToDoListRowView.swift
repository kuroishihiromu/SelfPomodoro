//
//  ToDoListRowView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/20.
//

import SwiftUI
import ComposableArchitecture

struct ToDoListRow: View {
    let store: StoreOf<ToDoListRowFeature>
    let width: CGFloat
    let height: CGFloat
    private let bgColor = ColorTheme.white
    private let fontColor = ColorTheme.black
    private let checkedColor = ColorTheme.navy
    private let uncheckedColor = ColorTheme.darkGray

    var body: some View {
        WithViewStore(store, observe: {$0}){ viewStore in
            HStack {
                Button(action: {
                    viewStore.send(.toggleCompleted)
                }) {
                    Image(systemName: viewStore.isCompleted ? "checkmark.square.fill" : "square")
                        .foregroundColor(viewStore.isCompleted ? checkedColor : uncheckedColor)
                        .font(.system(size: 16))
                }

                Text(viewStore.title)
                    .foregroundColor(fontColor)
                    .font(.system(size: 14))
            }
            .padding()
            .frame(maxWidth: width, maxHeight: height, alignment: .leading )
            .background(bgColor)
            .cornerRadius(8)
        }
    }
}

//#Preview {
//    PreviewTodoListRow(
//        store
//    )
//}
