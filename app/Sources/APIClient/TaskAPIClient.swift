//
//  TaskAPIClient.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/05/18.
//

import Foundation
import Dependencies

struct TaskResult: Equatable, Identifiable, Codable {
    let id: UUID
    var detail: String
    var isCompleted: Bool
}

enum taskAPIError: Error, Equatable {
    case networkError
    case decodingError
    case unknown
}

struct TaskAPIClient {
    var fetchTasks: () async throws -> [TaskResult]                         // GET /tasks
    var addTask: (_ detail: String) async throws -> TaskResult            // POST /tasks
    var deleteTask: (_ id: UUID) async throws -> Void                // DELETE /tasks/{id}
    var toggleCompletion: (_ id: UUID) async throws -> TaskResult          // PATCH /tasks/{id}/toggle
    var editTask: (_ id: UUID, _ detail: String) async throws -> TaskResult // PATCH /tasks/{id}/edit
}

extension TaskAPIClient {
    static let live = TaskAPIClient(
        fetchTasks: {
            var request = URLRequest(url: URL(string: "http://localhost:8080/api/v1/tasks")!)
            request.httpMethod = "GET"
            request.setValue("Bearer dev-token", forHTTPHeaderField: "Authorization")
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")

            let (data, response) = try await URLSession.shared.data(for: request)

            if let httpResponse = response as? HTTPURLResponse {
                print("Status Code: \(httpResponse.statusCode)")
            }

            print("Data fetched from tasks → \(String(data: data, encoding: .utf8) ?? "Invalid UTF-8")")

            return try JSONDecoder().decode([TaskResult].self, from: data)
        },

        addTask: { detail in
            var request = URLRequest(url: URL(string: "http://localhost:8080/api/v1/tasks")!)
            request.httpMethod = "POST"
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")
            request.httpBody = try JSONEncoder().encode(["detail": detail])
            let (data, _) = try await URLSession.shared.data(for: request)
            return try JSONDecoder().decode(TaskResult.self, from: data)
        },

        deleteTask: { id in
            var request = URLRequest(url: URL(string: "http://localhost:8080/api/v1/tasks/\(id)")!)
            request.httpMethod = "DELETE"
            _ = try await URLSession.shared.data(for: request)
        },

        toggleCompletion: { id in
            var request = URLRequest(url: URL(string: "http://localhost:8080/api/v1/tasks/\(id)/toggle")!)
            request.httpMethod = "PATCH"
            _ = try await URLSession.shared.data(for: request)
            let (data, _) = try await URLSession.shared.data(for: request)
            return try JSONDecoder().decode(TaskResult.self, from: data)
        },

        editTask: { id, detail in
            var request = URLRequest(url: URL(string: "http://localhost:8080/api/v1/tasks/\(id)/edit")!)
            request.httpMethod = "PATCH"
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")
            request.httpBody = try JSONEncoder().encode(["detail": detail])
            let (data, _) = try await URLSession.shared.data(for: request)
            return try JSONDecoder().decode(TaskResult.self, from: data)
        }
    )
}

extension DependencyValues {
    var taskAPIClient: TaskAPIClient {
        get { self[TaskAPIClientKey.self] }
        set { self[TaskAPIClientKey.self] = newValue }
    }

    private enum TaskAPIClientKey: DependencyKey {
        static let liveValue = TaskAPIClient.live
    }
}
