//
//  ErrorMessageHelper.swift
//  SelfPomodoro
//
//  Created by Claude on 2025/08/26.
//

import Foundation
import Amplify

struct ErrorMessageHelper {
    
    /// CognitoのAuthErrorを日本語化する
    static func localizedAuthError(_ error: Error) -> String {
        // AuthError以外はそのまま返す
        guard let authError = error as? AuthError else {
            return error.localizedDescription
        }
        
        // AuthErrorの詳細を解析
        switch authError {
        case .service(_, _, let underlyingError):
            return parseUnderlyingError(underlyingError)
        case .validation(_, _, _, _):
            return L10n.Error.authentication
        case .notAuthorized:
            return L10n.Error.authentication
        default:
            return L10n.Error.unknown
        }
    }
    
    /// 内部エラーを解析して適切な日本語メッセージを返す
    private static func parseUnderlyingError(_ underlyingError: Error?) -> String {
        guard let underlyingError = underlyingError else {
            return L10n.Error.unknown
        }
        
        let errorString = underlyingError.localizedDescription.lowercased()
        
        // パスワード関連エラー
        if errorString.contains("password not long enough") ||
           errorString.contains("password did not conform with policy") {
            return "パスワードは8文字以上で、英大文字・英小文字・数字を少なくとも1つ含めてください"
        }
        
        // メール関連エラー
        if errorString.contains("invalid email") || 
           errorString.contains("malformed email") {
            return "メールアドレスの形式が正しくありません"
        }
        
        if errorString.contains("email already exists") || 
           errorString.contains("usernameexistsexception") {
            return "このメールアドレスは既に使用されています"
        }
        
        // サインイン関連エラー
        if errorString.contains("incorrect username or password") ||
           errorString.contains("notauthorizedexception") {
            return "メールアドレスまたはパスワードが正しくありません"
        }
        
        if errorString.contains("user is not confirmed") ||
           errorString.contains("usernotconfirmedexception") {
            return "メールアドレスの確認が完了していません。確認メールをご確認ください"
        }
        
        if errorString.contains("user does not exist") ||
           errorString.contains("usernotfoundexception") {
            return "このメールアドレスで登録されたアカウントは見つかりません"
        }
        
        // 確認コード関連エラー
        if errorString.contains("invalid verification code") ||
           errorString.contains("codechallengeexception") {
            return "確認コードが正しくありません"
        }
        
        if errorString.contains("confirmation code expired") ||
           errorString.contains("expiredcodeexception") {
            return "確認コードの有効期限が切れています。新しい確認コードをリクエストしてください"
        }
        
        // レート制限エラー
        if errorString.contains("too many requests") ||
           errorString.contains("throttlingexception") {
            return "リクエストが多すぎます。しばらく時間をおいてから再試行してください"
        }
        
        // ネットワーク関連エラー
        if errorString.contains("network") || 
           errorString.contains("connection") {
            return L10n.Error.network
        }
        
        // その他は元のメッセージを返す
        return underlyingError.localizedDescription
    }
}
