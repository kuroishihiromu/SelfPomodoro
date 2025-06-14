// internal/handler/error_mapper.go
package handler

import (
	"errors"
	"net/http"

	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
)

// HTTPErrorResponse - HTTPレスポンス用エラー
type HTTPErrorResponse struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

// MapErrorToHTTP - エラーをHTTPレスポンスにマッピング
func MapErrorToHTTP(err error) *HTTPErrorResponse {
	// Domain Errors - ビジネスエラー（基本エラーから判定）
	switch {
	// User関連
	case errors.Is(err, appErrors.ErrUserNotFound):
		return &HTTPErrorResponse{
			Code:       "USER_NOT_FOUND",
			Message:    "ユーザーが見つかりません",
			StatusCode: http.StatusNotFound,
		}
	case errors.Is(err, appErrors.ErrUserEmailDuplicate):
		return &HTTPErrorResponse{
			Code:       "EMAIL_ALREADY_EXISTS",
			Message:    "メールアドレスが既に使用されています",
			StatusCode: http.StatusConflict,
		}
	case errors.Is(err, appErrors.ErrUserInvalidData):
		return &HTTPErrorResponse{
			Code:       "INVALID_USER_DATA",
			Message:    "無効なユーザーデータです",
			StatusCode: http.StatusBadRequest,
		}

	// Task関連
	case errors.Is(err, appErrors.ErrTaskNotFound):
		return &HTTPErrorResponse{
			Code:       "TASK_NOT_FOUND",
			Message:    "タスクが見つかりません",
			StatusCode: http.StatusNotFound,
		}
	case errors.Is(err, appErrors.ErrTaskAlreadyDone):
		return &HTTPErrorResponse{
			Code:       "TASK_ALREADY_DONE",
			Message:    "タスクは既に完了しています",
			StatusCode: http.StatusConflict,
		}
	case errors.Is(err, appErrors.ErrTaskInProgress):
		return &HTTPErrorResponse{
			Code:       "TASK_IN_PROGRESS",
			Message:    "既に進行中のタスクが存在します",
			StatusCode: http.StatusConflict,
		}

	// Session関連
	case errors.Is(err, appErrors.ErrSessionNotFound):
		return &HTTPErrorResponse{
			Code:       "SESSION_NOT_FOUND",
			Message:    "セッションが見つかりません",
			StatusCode: http.StatusNotFound,
		}
	case errors.Is(err, appErrors.ErrSessionInProgress):
		return &HTTPErrorResponse{
			Code:       "SESSION_IN_PROGRESS",
			Message:    "既に進行中のセッションが存在します",
			StatusCode: http.StatusConflict,
		}
	case errors.Is(err, appErrors.ErrSessionAlreadyEnded):
		return &HTTPErrorResponse{
			Code:       "SESSION_ALREADY_ENDED",
			Message:    "セッションは既に終了しています",
			StatusCode: http.StatusConflict,
		}

	// Round関連
	case errors.Is(err, appErrors.ErrRoundNotFound):
		return &HTTPErrorResponse{
			Code:       "ROUND_NOT_FOUND",
			Message:    "ラウンドが見つかりません",
			StatusCode: http.StatusNotFound,
		}
	case errors.Is(err, appErrors.ErrRoundAlreadyEnded):
		return &HTTPErrorResponse{
			Code:       "ROUND_ALREADY_ENDED",
			Message:    "ラウンドは既に終了しています",
			StatusCode: http.StatusConflict,
		}
	case errors.Is(err, appErrors.ErrRoundInProgress):
		return &HTTPErrorResponse{
			Code:       "ROUND_IN_PROGRESS",
			Message:    "既に進行中のラウンドが存在します",
			StatusCode: http.StatusConflict,
		}
	case errors.Is(err, appErrors.ErrNoRoundsInSession):
		return &HTTPErrorResponse{
			Code:       "NO_ROUNDS_IN_SESSION",
			Message:    "セッションにラウンドが存在しません",
			StatusCode: http.StatusNotFound,
		}

	// UserConfig関連
	case errors.Is(err, appErrors.ErrUserConfigNotFound):
		return &HTTPErrorResponse{
			Code:       "USER_CONFIG_NOT_FOUND",
			Message:    "ユーザー設定が見つかりません",
			StatusCode: http.StatusNotFound,
		}
	case errors.Is(err, appErrors.ErrUserConfigInvalid):
		return &HTTPErrorResponse{
			Code:       "INVALID_USER_CONFIG",
			Message:    "無効なユーザー設定です",
			StatusCode: http.StatusBadRequest,
		}

	// Statistics関連
	case errors.Is(err, appErrors.ErrStatisticsNotFound):
		return &HTTPErrorResponse{
			Code:       "STATISTICS_NOT_FOUND",
			Message:    "統計データが見つかりません",
			StatusCode: http.StatusNotFound,
		}

	// Authentication関連（ビジネス観点）
	case errors.Is(err, appErrors.ErrUnauthorized):
		return &HTTPErrorResponse{
			Code:       "UNAUTHORIZED",
			Message:    "認証が必要です",
			StatusCode: http.StatusUnauthorized,
		}
	case errors.Is(err, appErrors.ErrAccessDenied):
		return &HTTPErrorResponse{
			Code:       "ACCESS_DENIED",
			Message:    "アクセス権限がありません",
			StatusCode: http.StatusForbidden,
		}
	case errors.Is(err, appErrors.ErrTokenExpired):
		return &HTTPErrorResponse{
			Code:       "TOKEN_EXPIRED",
			Message:    "トークンの有効期限が切れています",
			StatusCode: http.StatusUnauthorized,
		}
	case errors.Is(err, appErrors.ErrInvalidToken):
		return &HTTPErrorResponse{
			Code:       "INVALID_TOKEN",
			Message:    "無効なトークンです",
			StatusCode: http.StatusUnauthorized,
		}
	case errors.Is(err, appErrors.ErrTokenNotFound):
		return &HTTPErrorResponse{
			Code:       "TOKEN_NOT_FOUND",
			Message:    "Authorizationヘッダーが見つかりません",
			StatusCode: http.StatusUnauthorized,
		}
	case errors.Is(err, appErrors.ErrInvalidFormat):
		return &HTTPErrorResponse{
			Code:       "INVALID_FORMAT",
			Message:    "リクエスト形式が無効です",
			StatusCode: http.StatusBadRequest,
		}

	// Validation関連
	case errors.Is(err, appErrors.ErrValidationFailed):
		return &HTTPErrorResponse{
			Code:       "VALIDATION_ERROR",
			Message:    "入力値が無効です",
			StatusCode: http.StatusBadRequest,
		}
	case errors.Is(err, appErrors.ErrRequiredField):
		return &HTTPErrorResponse{
			Code:       "REQUIRED_FIELD_MISSING",
			Message:    "必須項目が不足しています",
			StatusCode: http.StatusBadRequest,
		}
	case errors.Is(err, appErrors.ErrInvalidValue):
		return &HTTPErrorResponse{
			Code:       "INVALID_VALUE",
			Message:    "無効な値です",
			StatusCode: http.StatusBadRequest,
		}
	case errors.Is(err, appErrors.ErrInvalidEmail):
		return &HTTPErrorResponse{
			Code:       "INVALID_EMAIL",
			Message:    "無効なメールアドレスです",
			StatusCode: http.StatusBadRequest,
		}

	// Business Rule関連
	case errors.Is(err, appErrors.ErrBusinessRuleViolation):
		return &HTTPErrorResponse{
			Code:       "BUSINESS_RULE_VIOLATION",
			Message:    "ビジネスルール違反です",
			StatusCode: http.StatusConflict,
		}
	case errors.Is(err, appErrors.ErrConflict):
		return &HTTPErrorResponse{
			Code:       "CONFLICT",
			Message:    "リソースが既に存在します",
			StatusCode: http.StatusConflict,
		}
	case errors.Is(err, appErrors.ErrNotFound):
		return &HTTPErrorResponse{
			Code:       "NOT_FOUND",
			Message:    "リソースが見つかりません",
			StatusCode: http.StatusNotFound,
		}
	case errors.Is(err, appErrors.ErrBadRequest):
		return &HTTPErrorResponse{
			Code:       "BAD_REQUEST",
			Message:    "不正なリクエストです",
			StatusCode: http.StatusBadRequest,
		}
	case errors.Is(err, appErrors.ErrForbidden):
		return &HTTPErrorResponse{
			Code:       "FORBIDDEN",
			Message:    "アクセス権限がありません",
			StatusCode: http.StatusForbidden,
		}
	}

	// 構造化ドメインエラー
	if domainErr, ok := err.(*appErrors.DomainError); ok {
		return mapDomainErrorToHTTP(domainErr)
	}

	// Infrastructure Errors - 技術エラー（基本的に500）
	if appErrors.IsInfrastructureError(err) {
		return &HTTPErrorResponse{
			Code:       "INFRASTRUCTURE_ERROR",
			Message:    "システムエラーが発生しました",
			StatusCode: http.StatusInternalServerError,
		}
	}

	// Unknown Error
	return &HTTPErrorResponse{
		Code:       "INTERNAL_ERROR",
		Message:    "予期しないエラーが発生しました",
		StatusCode: http.StatusInternalServerError,
	}
}

// mapDomainErrorToHTTP - 構造化ドメインエラーをHTTPレスポンスにマッピング
func mapDomainErrorToHTTP(err *appErrors.DomainError) *HTTPErrorResponse {
	statusCode := http.StatusInternalServerError

	switch err.Code {
	case "USER_NOT_FOUND", "TASK_NOT_FOUND", "SESSION_NOT_FOUND", "ROUND_NOT_FOUND", "USER_CONFIG_NOT_FOUND", "NOT_FOUND":
		statusCode = http.StatusNotFound
	case "VALIDATION_ERROR", "BAD_REQUEST", "INVALID_USER_DATA", "INVALID_USER_CONFIG", "INVALID_FORMAT", "INVALID_VALUE", "INVALID_EMAIL", "REQUIRED_FIELD_MISSING":
		statusCode = http.StatusBadRequest
	case "UNAUTHORIZED", "TOKEN_EXPIRED", "INVALID_TOKEN", "TOKEN_NOT_FOUND":
		statusCode = http.StatusUnauthorized
	case "ACCESS_DENIED", "FORBIDDEN":
		statusCode = http.StatusForbidden
	case "BUSINESS_RULE_VIOLATION", "CONFLICT", "SESSION_IN_PROGRESS", "ROUND_IN_PROGRESS", "TASK_IN_PROGRESS", "SESSION_ALREADY_ENDED", "ROUND_ALREADY_ENDED", "TASK_ALREADY_DONE", "EMAIL_ALREADY_EXISTS":
		statusCode = http.StatusConflict
	case "USER_CREATION_FAILED", "USER_UPDATE_FAILED", "USER_DELETE_FAILED", "USER_CONFIG_CREATE_FAILED", "USER_CONFIG_UPDATE_FAILED", "USER_CONFIG_DELETE_FAILED":
		statusCode = http.StatusInternalServerError
	}

	return &HTTPErrorResponse{
		Code:       err.Code,
		Message:    err.Message,
		StatusCode: statusCode,
	}
}

// =================================
// ヘルパー関数
// =================================

// IsClientError - クライアントエラー（4xx）かどうかを判定
func IsClientError(err error) bool {
	httpErr := MapErrorToHTTP(err)
	return httpErr.StatusCode >= 400 && httpErr.StatusCode < 500
}

// IsServerError - サーバーエラー（5xx）かどうかを判定
func IsServerError(err error) bool {
	httpErr := MapErrorToHTTP(err)
	return httpErr.StatusCode >= 500
}

// IsNotFoundError - 404エラーかどうかを判定
func IsNotFoundError(err error) bool {
	httpErr := MapErrorToHTTP(err)
	return httpErr.StatusCode == http.StatusNotFound
}

// IsUnauthorizedError - 401エラーかどうかを判定
func IsUnauthorizedError(err error) bool {
	httpErr := MapErrorToHTTP(err)
	return httpErr.StatusCode == http.StatusUnauthorized
}

// IsForbiddenError - 403エラーかどうかを判定
func IsForbiddenError(err error) bool {
	httpErr := MapErrorToHTTP(err)
	return httpErr.StatusCode == http.StatusForbidden
}

// IsBadRequestError - 400エラーかどうかを判定
func IsBadRequestError(err error) bool {
	httpErr := MapErrorToHTTP(err)
	return httpErr.StatusCode == http.StatusBadRequest
}

// IsConflictError - 409エラーかどうかを判定
func IsConflictError(err error) bool {
	httpErr := MapErrorToHTTP(err)
	return httpErr.StatusCode == http.StatusConflict
}

// GetHTTPStatusCode - エラーからHTTPステータスコードを取得
func GetHTTPStatusCode(err error) int {
	return MapErrorToHTTP(err).StatusCode
}

// GetErrorCode - エラーからエラーコードを取得
func GetErrorCode(err error) string {
	return MapErrorToHTTP(err).Code
}

// GetErrorMessage - エラーからエラーメッセージを取得
func GetErrorMessage(err error) string {
	return MapErrorToHTTP(err).Message
}
