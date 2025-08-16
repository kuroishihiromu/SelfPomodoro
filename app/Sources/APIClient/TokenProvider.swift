//
//  TokenProvider.swift
//  SelfPomodoro
//
//  Created by Claude on 2025/01/16.
//

import Dependencies
import Foundation

/// トークンプロバイダー - 各APIClientでAuthorizationヘッダーを提供
struct TokenProvider {
    var getAuthToken: () -> String?
}

extension TokenProvider {
    static let live = TokenProvider(
        getAuthToken: {
            return TokenStorage.shared.currentToken?.idToken
        }
    )

    static let preview = TokenProvider(
        getAuthToken: { "preview-mock-token" }
    )

    static let test = TokenProvider(
        getAuthToken: { nil }
    )
}

extension DependencyValues {
    var tokenProvider: TokenProvider {
        get { self[TokenProviderKey.self] }
        set { self[TokenProviderKey.self] = newValue }
    }

    private enum TokenProviderKey: DependencyKey {
        static let liveValue = TokenProvider.live
        static let previewValue = TokenProvider.preview
        static let testValue = TokenProvider.test
    }
}

/// トークンストレージ - グローバルなトークン管理
class TokenStorage {
    static let shared = TokenStorage()

    private var _currentToken: AuthTokens?

    private init() {}

    var currentToken: AuthTokens? {
        return _currentToken
    }

    func setToken(_ token: AuthTokens?) {
        _currentToken = token
    }

    func clearToken() {
        _currentToken = nil
    }
}
