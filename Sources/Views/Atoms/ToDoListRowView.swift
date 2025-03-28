//
//  ToDoListRowView.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/03/20.
//

import SwiftUI

struct TodoListRow: View {
    let title: String
    let width: CGFloat
    let height: CGFloat
    @Binding var isCompleted: Bool

    private let bgColor = ColorTheme.white
    private let fontColor = ColorTheme.black
    private let checkedColor = ColorTheme.navy
    private let uncheckedColor = ColorTheme.darkGray

    init(
        title: String,
        width: CGFloat,
        height: CGFloat,
        isCompleted: Binding<Bool>
    ) {
        self.title = title
        self.width = width
        self.height = height
        self._isCompleted = isCompleted
    }

    var body: some View {
        HStack(spacing: 5) {
            Button(action: {
                isCompleted.toggle()
            }) {
                Image(systemName: isCompleted ? "checkmark.square.fill" : "square")
                    .foregroundColor(isCompleted ? checkedColor : uncheckedColor)
                    .font(.system(size: 16))
            }

            Text(title)
                .foregroundColor(fontColor)
                .font(.system(size: 14))
        }
        .padding()
        .frame(width: width, height: height, alignment: .leading)
        .background(bgColor)
    }
}

#Preview {
    PreviewTodoListRow()
}

struct PreviewTodoListRow: View {
    @State var isChecked = false

    var body: some View {
        VStack {
            TodoListRow(
                title: "Task 1: Review design specifications",
                width: 300,
                height: 30,
                isCompleted: $isChecked
            )
        }
    }
}
