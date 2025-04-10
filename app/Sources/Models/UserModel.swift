//
//  UserModel.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/02/27.
//

import Foundation

struct Profile: Decodable {
    let username: String?
    let fullName: String?
    let website: String?

    enum CodingKeys: String, CodingKey {
        case username
        case fullName = "full_name"
        case website
    }
}

struct UpdateProfileParams: Encodable {
    let username: String
    let fullName: String
    let website: String

    enum CodingKeys: String, CodingKey {
        case username
        case fullName = "full_name"
        case website
    }
}
