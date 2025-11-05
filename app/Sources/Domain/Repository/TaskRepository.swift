//
//  TaskRepository.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/05.
//

import Foundation

protocol TaskRepository {
    func fetchTasks(for identifier: String) async throws -> [TodoTask]
    func createTask(detail: String, for identifier: String) async throws -> TodoTask
    func toggleTaskCompletion(for id: UUID) async throws -> TodoTask
    func deleteTask(with id: UUID) async throws
}
