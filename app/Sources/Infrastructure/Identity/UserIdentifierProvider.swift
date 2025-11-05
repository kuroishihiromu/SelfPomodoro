//
//  UserIdentifierProvider.swift
//  SelfPomodoro
//
//  Created by tsunakit99 on 2025/11/05.
//

import Foundation
import UIKit

enum UserIdentifierProvider {
    private static let storageKey = "com.selfpomodoro.userIdentifier"

    static func resolve() -> String {
        if let stored = UserDefaults.standard.string(forKey: storageKey) {
            return stored
        }

        if let vendorId = UIDevice.current.identifierForVendor?.uuidString {
            store(vendorId)
            return vendorId
        }

        let fresh = UUID().uuidString
        store(fresh)
        return fresh
    }

    private static func store(_ value: String) {
        UserDefaults.standard.set(value, forKey: storageKey)
    }
}
