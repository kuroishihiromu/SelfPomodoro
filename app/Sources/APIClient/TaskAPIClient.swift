//
//  TaskAPIClient.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/05/18.
//

import Foundation
import Dependencies
import Amplify

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
            let request = RESTRequest(
                apiName: "selfpomodoro",
                path: "/api/v1/tasks",
                headers: [
                    "Content-Type": "application/json"
                ]
            )

            let data = try await Amplify.API.get(request: request)
            print("Fetched tasks → \(String(data: data, encoding: .utf8) ?? "Invalid UTF-8")")

            struct TaskListResponse: Decodable {
                let tasks: [TaskResult]
            }

            return try APIFormatters.jsonDecoder.decode(TaskListResponse.self, from: data).tasks
        },
        
        addTask: { detail in
            let body = try JSONEncoder().encode(["detail": detail])
            let request = RESTRequest(
                path: "/api/v1/tasks",
                headers: [
                    "Content-Type": "application/json"
                ],
                body: body
            )

            let data = try await Amplify.API.post(request: request)
            print("Add task response → \(String(data: data, encoding: .utf8) ?? "Invalid UTF-8")")
            return try APIFormatters.jsonDecoder.decode(TaskResult.self, from: data)
        },

        deleteTask: { id in
            let request = RESTRequest(
                path: "/api/v1/tasks/\(id.uuidString)",
                headers: [
                    "Content-Type": "application/json"
                ]
            )

            _ = try await Amplify.API.delete(request: request)
        },

        toggleCompletion: { id in
            let request = RESTRequest(
                apiName: "selfpomodoro",
                path: "/api/v1/tasks/\(id)/toggle",
                headers: ["Content-Type": "application/json"],
                body: nil
            )

            let data = try await Amplify.API.patch(request: request)

            print("Toggle task response → \(String(data: data, encoding: .utf8) ?? "Invalid UTF-8")")

            return try APIFormatters.jsonDecoder.decode(TaskResult.self, from: data)
        },
        
        editTask: { id, detail in
            var request = URLRequest(url: URL(string: "http://localhost:8080/api/v1/tasks/\(id)/edit")!)
            request.httpMethod = "PATCH"
            request.setValue("Bearer dev-token", forHTTPHeaderField: "Authorization")
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")
            request.httpBody = try JSONEncoder().encode(["detail": detail])

            let (data, _) = try await URLSession.shared.data(for: request)
            print("Edit task response → \(String(data: data, encoding: .utf8) ?? "Invalid UTF-8")")
            return try APIFormatters.jsonDecoder.decode(TaskResult.self, from: data)
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
