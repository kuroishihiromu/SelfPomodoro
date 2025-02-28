//
//  EnvironmentKey.swift
//  SelfPomodoro
//
//  Created by 黒石陽夢 on 2025/02/27.
//

import SwiftUI
struct AuthStatusKey : EnvironmentKey {
    static var defaultValue: Bool = false
}

extension EnvironmentValues {
    var isAuthenticated : Bool {
        get {
            self[AuthStatusKey.self]
        } set {
            self[AuthStatusKey.self] = newValue
        }
    }
}
