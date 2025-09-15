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
}

extension TaskAPIClient {
    static let live = TaskAPIClient(
        fetchTasks: {
            do {
                let session = try await Amplify.Auth.fetchAuthSession()
                print("👤 Auth session (fetchTasks) isSignedIn=\(session.isSignedIn)")
            } catch {
                print("👤 Auth session (fetchTasks) fetch failed: \(error)")
            }
            print("➡️ GET /dev/api/v1/tasks")
            let request = RESTRequest(
                apiName: "selfpomodoro",
                path: "/dev/api/v1/tasks"
            )

            do {
                let data = try await Amplify.API.get(request: request)
                print("📦 fetchTasks bytes=\(data.count)")
                struct TaskListResponse: Decodable { let tasks: [TaskResult] }
                let tasks = try AppDecoder.default.decode(TaskListResponse.self, from: data).tasks
                print("✅ fetchTasks count=\(tasks.count)")
                return tasks
            } catch {
                print("❌ fetchTasks failed: \(error)")
                throw error
            }

        },
        
        addTask: { detail in
            do {
                let session = try await Amplify.Auth.fetchAuthSession()
                print("👤 Auth session (addTask) isSignedIn=\(session.isSignedIn)")
            } catch {
                print("👤 Auth session (addTask) fetch failed: \(error)")
            }
            print("➡️ POST /dev/api/v1/tasks")
            let body = try JSONEncoder().encode(["detail": detail])
            let request = RESTRequest(
                apiName: "selfpomodoro",
                path: "/dev/api/v1/tasks",
                body: body
            )

            do {
                let data = try await Amplify.API.post(request: request)
                print("📦 addTask bytes=\(data.count)")
                let task = try AppDecoder.default.decode(TaskResult.self, from: data)
                print("✅ addTask id=\(task.id)")
                return task
            } catch {
                print("❌ addTask failed: \(error)")
                throw error
            }
        },

        deleteTask: { id in
            do {
                let session = try await Amplify.Auth.fetchAuthSession()
                print("👤 Auth session (deleteTask) isSignedIn=\(session.isSignedIn)")
            } catch {
                print("👤 Auth session (deleteTask) fetch failed: \(error)")
            }
            print("➡️ DELETE /dev/api/v1/tasks/\(id.uuidString)")
            let request = RESTRequest(
                apiName: "selfpomodoro",
                path: "/dev/api/v1/tasks/\(id.uuidString)"
            )

            do {
                _ = try await Amplify.API.delete(request: request)
                print("✅ deleteTask id=\(id)")
            } catch {
                print("❌ deleteTask failed: \(error)")
                throw error
            }
        },

        toggleCompletion: { id in
            do {
                let session = try await Amplify.Auth.fetchAuthSession()
                print("👤 Auth session (toggleCompletion) isSignedIn=\(session.isSignedIn)")
            } catch {
                print("👤 Auth session (toggleCompletion) fetch failed: \(error)")
            }
            print("➡️ PATCH /dev/api/v1/tasks/\(id)/toggle")
            let request = RESTRequest(
                apiName: "selfpomodoro",
                path: "/dev/api/v1/tasks/\(id)/toggle"
            )

            do {
                let data = try await Amplify.API.patch(request: request)
                print("📦 toggleCompletion bytes=\(data.count)")
                let task = try AppDecoder.default.decode(TaskResult.self, from: data)
                print("✅ toggleCompletion id=\(task.id) completed=\(task.isCompleted)")
                return task
            } catch {
                print("❌ toggleCompletion failed: \(error)")
                throw error
            }
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
