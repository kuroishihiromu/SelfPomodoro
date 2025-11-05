//
//  SwiftDataUserRepository.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/05.
//

import Foundation
import SwiftData

@MainActor
final class SwiftDataUserRepository: UserRepository {
    private let context: ModelContext

    init(context: ModelContext) {
       self.context = context
   }

    // MARK: - UserRepository

    func fetchUser(by identifier: String) async throws -> User? {
        guard let model = try fetchUserModel(by: identifier) else {
            return nil
        }
        return mapUser(model)
    }

    func createUser(identifier: String) async throws -> User {
        let now = Date()
        let model = UserModel(identifier: identifier, createdAt: now, updatedAt: now)
        context.insert(model)
        try context.save()
        return mapUser(model)
    }

    // MARK: - Helpers

    private func fetchUserModel(by identifier: String) throws -> UserModel? {
        let descriptor = FetchDescriptor<UserModel>(
            predicate: #Predicate { $0.identifier == identifier },
            fetchLimit: 1
        )
        return try context.fetch(descriptor).first
    }

    private func mapUser(_ model: UserModel) -> User {
        User(
            id: model.id,
            identifier: model.identifier,
            createdAt: model.createdAt,
            updatedAt: model.updatedAt
        )
    }
}
