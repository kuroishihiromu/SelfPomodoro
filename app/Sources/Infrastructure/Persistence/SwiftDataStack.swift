//
//  SwiftDataStack.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/05.
//

import Foundation
import SwiftData

enum SwiftDataStack {
    @MainActor
    static func makeContainer(inMemory: Bool = false) -> ModelContainer {
        let configuration = ModelConfiguration(isStoredInMemoryOnly: inMemory)
        do {
            return try ModelContainer(
                for: UserModel.self,
                TaskModel.self,
                configurations: configuration
            )
        } catch {
            fatalError("Failed to create SwiftData container: \(error.localizedDescription)")
        }
    }
}
