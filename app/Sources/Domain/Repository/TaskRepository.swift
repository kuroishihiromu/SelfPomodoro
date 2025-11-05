//
//  TaskRepository.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/05.
//

import Foundation

protocol TaskRepository {
    func fetchTasks(for identifier: String) async throws -> [Task]
    func createTask(detail: String, for identifier: String) async throws -> Task
    func toggleTaskCompletion(for id: UUID) async throws -> Task
    func deleteTask(with id: UUID) async throws
}
