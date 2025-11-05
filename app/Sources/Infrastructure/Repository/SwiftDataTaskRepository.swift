//
//  SwiftDataTaskRepository.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/05.
//

import Foundation
import SwiftData

@MainActor
final class SwiftDataTaskRepository: TaskRepository {
    private let context: ModelContext

    init(context: ModelContext) {
        self.context = context
    }

    // MARK: - TaskRepository

    func createTask(detail: String, for identifier: String) async throws -> TodoTask {
        let now = Date()
        let model = TaskModel(
            userIdentifier: identifier,
            detail: detail,
            createdAt: now,
            updatedAt: now,
            user: nil
        )
        context.insert(model)
        try context.save()
        print("📝 SwiftDataTaskRepository: created task id=\(model.id) detail=\(detail) userIdentifier=\(identifier)")
        return mapTask(model)
    }

    func fetchTasks(for identifier: String) async throws -> [TodoTask] {
        var descriptor = FetchDescriptor<TaskModel>(
            predicate: #Predicate { $0.userIdentifier == identifier },
            sortBy: [SortDescriptor(\.createdAt, order: .forward)]
        )
        let models = try context.fetch(descriptor)
        print("📦 SwiftDataTaskRepository: fetched \(models.count) tasks for identifier=\(identifier)")
        return models.map(mapTask)
    }

    func toggleTaskCompletion(for id: UUID) async throws -> TodoTask {
        guard let model = try fetchTaskModel(by: id) else {
            throw RepositoryError.recordNotFound
        }
        model.isCompleted.toggle()
        model.updatedAt = Date()
        try context.save()
        print("🔁 SwiftDataTaskRepository: toggled task id=\(id) isCompleted=\(model.isCompleted)")
        return mapTask(model)
    }

    func deleteTask(with id: UUID) async throws {
        guard let model = try fetchTaskModel(by: id) else { return }
        context.delete(model)
        try context.save()
        print("🗑️ SwiftDataTaskRepository: deleted task id=\(id)")
    }

    // MARK: - Helpers

    private func fetchTaskModel(by id: UUID) throws -> TaskModel? {
        var descriptor = FetchDescriptor<TaskModel>(
            predicate: #Predicate { $0.id == id }
        )
        descriptor.fetchLimit = 1
        return try context.fetch(descriptor).first
    }

    private func mapTask(_ model: TaskModel) -> TodoTask {
        TodoTask(
            id: model.id,
            detail: model.detail,
            isCompleted: model.isCompleted,
            createdAt: model.createdAt,
            updatedAt: model.updatedAt
        )
    }

    enum RepositoryError: Error {
        case recordNotFound
    }
}
