//
//  UserRepository.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/05.
//

import Foundation

protocol UserRepository {
    func fetchUser(by identifier: String) async throws -> User?
    func createUser(identifier: String) async throws -> User
}
