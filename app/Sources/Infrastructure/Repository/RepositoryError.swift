//
//  RepositoryError.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/07.
//

import Foundation

enum RepositoryError: Error {
    case userNotFound
    case recordNotFound
    case invalidIdentifier
}
