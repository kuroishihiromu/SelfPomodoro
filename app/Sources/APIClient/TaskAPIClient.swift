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
    var createdAt: Date?
    var updatedAt: Date?

    enum CodingKeys: String, CodingKey {
        case id
        case detail
        case isCompleted = "is_completed"
        case createdAt = "created_at"
        case updatedAt = "updated_at"
    }
}

enum taskAPIError: Error, Equatable {
    case networkError
    case decodingError
    case unknown
}

struct TaskAPIClient {
    var fetchTasks: () async throws -> [TaskResult]
    var addTask: (_ detail: String) async throws -> TaskResult
    var deleteTask: (_ id: UUID) async throws -> Void
    var toggleCompletion: (_ id: UUID) async throws -> TaskResult
    var editTask: (_ id: UUID, _ detail: String) async throws -> TaskResult
}

extension TaskAPIClient {
    static let live = TaskAPIClient(
        fetchTasks: {
            var request = URLRequest(url: URL(string: "http://localhost:8080/api/v1/tasks")!)
            request.httpMethod = "GET"
            request.setValue("Bearer dev-token", forHTTPHeaderField: "Authorization")
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")

            let (data, response) = try await URLSession.shared.data(for: request)
            print("Fetched tasks → \(String(data: data, encoding: .utf8) ?? "Invalid UTF-8")")

            let formatter = ISO8601DateFormatter()
            formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]

            let decoder = JSONDecoder()
            decoder.dateDecodingStrategy = .custom { decoder in
                let container = try decoder.singleValueContainer()
                let dateStr = try container.decode(String.self)
                guard let date = formatter.date(from: dateStr) else {
                    throw DecodingError.dataCorruptedError(
                        in: container,
                        debugDescription: "Invalid date format: \(dateStr)"
                    )
                }
                return date
            }

            struct TaskListResponse: Decodable {
                let tasks: [TaskResult]
            }

            return try decoder.decode(TaskListResponse.self, from: data).tasks
        },


        addTask: { detail in
            var request = URLRequest(url: URL(string: "http://localhost:8080/api/v1/tasks")!)
            request.httpMethod = "POST"
            request.setValue("Bearer dev-token", forHTTPHeaderField: "Authorization")
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")
            request.httpBody = try JSONEncoder().encode(["detail": detail])

            let (data, _) = try await URLSession.shared.data(for: request)
            print("Add task response → \(String(data: data, encoding: .utf8) ?? "Invalid UTF-8")")
            
            let decoder = JSONDecoder()

            let formatter = ISO8601DateFormatter()
            formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds, .withTimeZone] // 重要！

            decoder.dateDecodingStrategy = .custom { decoder in
                let container = try decoder.singleValueContainer()
                let dateStr = try container.decode(String.self)
                guard let date = formatter.date(from: dateStr) else {
                    throw DecodingError.dataCorruptedError(
                        in: container,
                        debugDescription: "Invalid date format: \(dateStr)"
                    )
                }
                return date
            }
            return try decoder.decode(TaskResult.self, from: data)
        },

        deleteTask: { id in
            var request = URLRequest(url: URL(string: "http://localhost:8080/api/v1/tasks/\(id)")!)
            request.httpMethod = "DELETE"
            request.setValue("Bearer dev-token", forHTTPHeaderField: "Authorization")
            _ = try await URLSession.shared.data(for: request)
        },

        toggleCompletion: { id in
            var request = URLRequest(url: URL(string: "http://localhost:8080/api/v1/tasks/\(id)/toggle")!)
            request.httpMethod = "PATCH"
            request.setValue("Bearer dev-token", forHTTPHeaderField: "Authorization")
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")

            let (data, _) = try await URLSession.shared.data(for: request)
            print("Toggle task response → \(String(data: data, encoding: .utf8) ?? "Invalid UTF-8")")

            let formatter = ISO8601DateFormatter()
            formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]

            let decoder = JSONDecoder()
            decoder.dateDecodingStrategy = .custom { decoder in
                let container = try decoder.singleValueContainer()
                let dateStr = try container.decode(String.self)
                guard let date = formatter.date(from: dateStr) else {
                    throw DecodingError.dataCorruptedError(
                        in: container,
                        debugDescription: "Invalid date format: \(dateStr)"
                    )
                }
                return date
            }

            return try decoder.decode(TaskResult.self, from: data)
        },

        editTask: { id, detail in
            var request = URLRequest(url: URL(string: "http://localhost:8080/api/v1/tasks/\(id)/edit")!)
            request.httpMethod = "PATCH"
            request.setValue("Bearer dev-token", forHTTPHeaderField: "Authorization")
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")
            request.httpBody = try JSONEncoder().encode(["detail": detail])

            let (data, _) = try await URLSession.shared.data(for: request)
            print("Edit task response → \(String(data: data, encoding: .utf8) ?? "Invalid UTF-8")")
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
