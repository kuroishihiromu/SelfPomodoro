package errors

import (
	"errors"
	"fmt"
)

// =================================
// 基本ビジネスエラー（HTTPステータスなし）
// =================================

var (
	// User関連
	ErrUserNotFound       = errors.New("user not found")
	ErrUserEmailDuplicate = errors.New("email already exists")
	ErrUserInvalidData    = errors.New("invalid user data")
	ErrUserCreationFailed = errors.New("user creation failed")
	ErrUserUpdateFailed   = errors.New("user update failed")
	ErrUserDeleteFailed   = errors.New("user delete failed")

	// Task関連
	ErrTaskNotFound     = errors.New("task not found")
	ErrTaskAlreadyDone  = errors.New("task already completed")
	ErrTaskInvalidState = errors.New("invalid task state")
	ErrTaskInProgress   = errors.New("task in progress")
	ErrTaskInvalidType  = errors.New("invalid task type")

	// Session関連
	ErrSessionNotFound     = errors.New("session not found")
	ErrSessionInProgress   = errors.New("session already in progress")
	ErrSessionAlreadyEnded = errors.New("session already ended")
	ErrSessionInvalidState = errors.New("invalid session state")

	// Round関連
	ErrRoundNotFound     = errors.New("round not found")
	ErrRoundAlreadyEnded = errors.New("round already ended")
	ErrRoundInProgress   = errors.New("round in progress")
	ErrRoundInvalidState = errors.New("invalid round state")
	ErrRoundInvalidType  = errors.New("invalid round type")
	ErrNoRoundsInSession = errors.New("no rounds in session")

	// UserConfig関連
	ErrUserConfigNotFound     = errors.New("user config not found")
	ErrUserConfigInvalid      = errors.New("invalid user config")
	ErrUserConfigCreateFailed = errors.New("user config creation failed")
	ErrUserConfigUpdateFailed = errors.New("user config update failed")
	ErrUserConfigDeleteFailed = errors.New("user config delete failed")

	// Statistics関連
	ErrStatisticsNotFound = errors.New("statistics not found")
	ErrStatisticsInvalid  = errors.New("invalid statistics data")

	// Authentication関連（ビジネス観点）
	ErrUnauthorized     = errors.New("authentication required")
	ErrAccessDenied     = errors.New("access denied")
	ErrTokenExpired     = errors.New("token expired")
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenNotFound    = errors.New("token not found")
	ErrInvalidFormat    = errors.New("invalid format")
	ErrInvalidIssuer    = errors.New("invalid issuer")
	ErrInvalidAudience  = errors.New("invalid audience")
	ErrInvalidSignature = errors.New("invalid signature")
	ErrMissingSubject   = errors.New("missing subject")
	ErrInvalidSubject   = errors.New("invalid subject")
	ErrMissingTokenUse  = errors.New("missing token use")
	ErrInvalidTokenUse  = errors.New("invalid token use")

	// Validation関連
	ErrValidationFailed = errors.New("validation failed")
	ErrRequiredField    = errors.New("required field missing")
	ErrInvalidValue     = errors.New("invalid value")
	ErrInvalidEmail     = errors.New("invalid email")

	// Business Rule関連
	ErrBusinessRuleViolation = errors.New("business rule violation")
	ErrConflict              = errors.New("resource conflict")
	ErrNotFound              = errors.New("resource not found")
	ErrBadRequest            = errors.New("bad request")
	ErrForbidden             = errors.New("forbidden")
	ErrInternalError         = errors.New("internal error")
)

// =================================
// 構造化ドメインエラー
// =================================

type DomainError struct {
	Code    string
	Message string
	Cause   error
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Cause
}

// =================================
// ドメインエラー作成ヘルパー関数
// =================================

func NewDomainError(code, message string, cause error) *DomainError {
	return &DomainError{Code: code, Message: message, Cause: cause}
}

// User関連エラー
func NewUserNotFoundError() error {
	return NewDomainError("USER_NOT_FOUND", "ユーザーが見つかりません", ErrUserNotFound)
}

func NewEmailAlreadyExistsError() error {
	return NewDomainError("EMAIL_ALREADY_EXISTS", "メールアドレスが既に使用されています", ErrUserEmailDuplicate)
}

func NewUserCreationFailedError() error {
	return NewDomainError("USER_CREATION_FAILED", "ユーザーの作成に失敗しました", ErrUserCreationFailed)
}

func NewUserUpdateFailedError() error {
	return NewDomainError("USER_UPDATE_FAILED", "ユーザーの更新に失敗しました", ErrUserUpdateFailed)
}

func NewUserDeleteFailedError() error {
	return NewDomainError("USER_DELETE_FAILED", "ユーザーの削除に失敗しました", ErrUserDeleteFailed)
}

func NewInvalidUserDataError() error {
	return NewDomainError("INVALID_USER_DATA", "無効なユーザーデータです", ErrUserInvalidData)
}

// Task関連エラー
func NewTaskNotFoundError() error {
	return NewDomainError("TASK_NOT_FOUND", "タスクが見つかりません", ErrTaskNotFound)
}

func NewTaskAlreadyDoneError() error {
	return NewDomainError("TASK_ALREADY_DONE", "タスクは既に完了しています", ErrTaskAlreadyDone)
}

func NewInvalidTaskStateError() error {
	return NewDomainError("INVALID_TASK_STATE", "タスクの状態が無効です", ErrTaskInvalidState)
}

func NewTaskInProgressError() error {
	return NewDomainError("TASK_IN_PROGRESS", "既に進行中のタスクが存在します", ErrTaskInProgress)
}

func NewInvalidTaskTypeError() error {
	return NewDomainError("INVALID_TASK_TYPE", "無効なタスクタイプです", ErrTaskInvalidType)
}

// Session関連エラー
func NewSessionNotFoundError() error {
	return NewDomainError("SESSION_NOT_FOUND", "セッションが見つかりません", ErrSessionNotFound)
}

func NewSessionInProgressError() error {
	return NewDomainError("SESSION_IN_PROGRESS", "既に進行中のセッションが存在します", ErrSessionInProgress)
}

func NewSessionAlreadyEndedError() error {
	return NewDomainError("SESSION_ALREADY_ENDED", "セッションは既に終了しています", ErrSessionAlreadyEnded)
}

func NewInvalidSessionStateError() error {
	return NewDomainError("INVALID_SESSION_STATE", "セッションの状態が無効です", ErrSessionInvalidState)
}

// Round関連エラー
func NewRoundNotFoundError() error {
	return NewDomainError("ROUND_NOT_FOUND", "ラウンドが見つかりません", ErrRoundNotFound)
}

func NewRoundAlreadyEndedError() error {
	return NewDomainError("ROUND_ALREADY_ENDED", "ラウンドは既に終了しています", ErrRoundAlreadyEnded)
}

func NewRoundInProgressError() error {
	return NewDomainError("ROUND_IN_PROGRESS", "既に進行中のラウンドが存在します", ErrRoundInProgress)
}

func NewInvalidRoundStateError() error {
	return NewDomainError("INVALID_ROUND_STATE", "ラウンドの状態が無効です", ErrRoundInvalidState)
}

func NewInvalidRoundTypeError() error {
	return NewDomainError("INVALID_ROUND_TYPE", "無効なラウンドタイプです", ErrRoundInvalidType)
}

// UserConfig関連エラー
func NewUserConfigNotFoundError() error {
	return NewDomainError("USER_CONFIG_NOT_FOUND", "ユーザー設定が見つかりません", ErrUserConfigNotFound)
}

func NewUserConfigCreateFailedError() error {
	return NewDomainError("USER_CONFIG_CREATE_FAILED", "ユーザー設定の作成に失敗しました", ErrUserConfigCreateFailed)
}

func NewUserConfigUpdateFailedError() error {
	return NewDomainError("USER_CONFIG_UPDATE_FAILED", "ユーザー設定の更新に失敗しました", ErrUserConfigUpdateFailed)
}

func NewUserConfigDeleteFailedError() error {
	return NewDomainError("USER_CONFIG_DELETE_FAILED", "ユーザー設定の削除に失敗しました", ErrUserConfigDeleteFailed)
}

func NewInvalidUserConfigError() error {
	return NewDomainError("INVALID_USER_CONFIG", "無効なユーザー設定です", ErrUserConfigInvalid)
}

// Authentication関連エラー
func NewUnauthorizedError(message string) error {
	if message == "" {
		message = "認証が必要です"
	}
	return NewDomainError("UNAUTHORIZED", message, ErrUnauthorized)
}

func NewTokenNotFoundError() error {
	return NewDomainError("TOKEN_NOT_FOUND", "Authorizationヘッダーが見つかりません", ErrTokenNotFound)
}

func NewInvalidTokenError() error {
	return NewDomainError("INVALID_TOKEN", "無効なJWTトークンです", ErrInvalidToken)
}

func NewTokenExpiredError() error {
	return NewDomainError("TOKEN_EXPIRED", "JWTトークンの有効期限が切れています", ErrTokenExpired)
}

func NewInvalidFormatError() error {
	return NewDomainError("INVALID_FORMAT", "Authorizationヘッダーの形式が無効です", ErrInvalidFormat)
}

func NewInvalidIssuerError() error {
	return NewDomainError("INVALID_ISSUER", "トークンの発行者が無効です", ErrInvalidIssuer)
}

func NewInvalidAudienceError() error {
	return NewDomainError("INVALID_AUDIENCE", "トークンのaudienceが無効です", ErrInvalidAudience)
}

func NewInvalidSignatureError() error {
	return NewDomainError("INVALID_SIGNATURE", "トークンの署名が無効です", ErrInvalidSignature)
}

func NewMissingSubjectError() error {
	return NewDomainError("MISSING_SUBJECT", "subjectクレームが見つかりません", ErrMissingSubject)
}

func NewInvalidSubjectError() error {
	return NewDomainError("INVALID_SUBJECT", "subjectクレームが無効です", ErrInvalidSubject)
}

func NewMissingTokenUseError() error {
	return NewDomainError("MISSING_TOKEN_USE", "token_useクレームが見つかりません", ErrMissingTokenUse)
}

func NewInvalidTokenUseError() error {
	return NewDomainError("INVALID_TOKEN_USE", "token_useクレームが無効です", ErrInvalidTokenUse)
}

// 汎用エラー
func NewValidationError(message string) error {
	if message == "" {
		message = "バリデーションエラーが発生しました"
	}
	return NewDomainError("VALIDATION_ERROR", message, ErrValidationFailed)
}

func NewBadRequestError(message string) error {
	if message == "" {
		message = "不正なリクエストです"
	}
	return NewDomainError("BAD_REQUEST", message, ErrBadRequest)
}

func NewNotFoundError(message string) error {
	if message == "" {
		message = "リソースが見つかりません"
	}
	return NewDomainError("NOT_FOUND", message, ErrNotFound)
}

func NewForbiddenError(message string) error {
	if message == "" {
		message = "アクセス権限がありません"
	}
	return NewDomainError("FORBIDDEN", message, ErrForbidden)
}

func NewConflictError(message string) error {
	if message == "" {
		message = "リソースが既に存在します"
	}
	return NewDomainError("CONFLICT", message, ErrConflict)
}

func NewInternalError(cause error) error {
	return NewDomainError("INTERNAL_ERROR", "内部サーバーエラーが発生しました", cause)
}

// =================================
// エラー判定関数
// =================================

func IsBusinessError(err error) bool {
	return errors.Is(err, ErrUserNotFound) ||
		errors.Is(err, ErrTaskNotFound) ||
		errors.Is(err, ErrSessionNotFound) ||
		errors.Is(err, ErrRoundNotFound) ||
		errors.Is(err, ErrUserConfigNotFound) ||
		errors.Is(err, ErrUnauthorized) ||
		errors.Is(err, ErrValidationFailed) ||
		errors.Is(err, ErrBusinessRuleViolation)
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrUserNotFound) ||
		errors.Is(err, ErrTaskNotFound) ||
		errors.Is(err, ErrSessionNotFound) ||
		errors.Is(err, ErrRoundNotFound) ||
		errors.Is(err, ErrUserConfigNotFound) ||
		errors.Is(err, ErrStatisticsNotFound)
}

func IsAuthError(err error) bool {
	return errors.Is(err, ErrUnauthorized) ||
		errors.Is(err, ErrAccessDenied) ||
		errors.Is(err, ErrTokenExpired) ||
		errors.Is(err, ErrInvalidToken) ||
		errors.Is(err, ErrTokenNotFound)
}

func IsValidationError(err error) bool {
	return errors.Is(err, ErrValidationFailed) ||
		errors.Is(err, ErrRequiredField) ||
		errors.Is(err, ErrInvalidValue) ||
		errors.Is(err, ErrInvalidEmail)
}

func IsConflictError(err error) bool {
	return errors.Is(err, ErrSessionInProgress) ||
		errors.Is(err, ErrRoundInProgress) ||
		errors.Is(err, ErrTaskInProgress) ||
		errors.Is(err, ErrConflict)
}

func IsAuthenticationError(err error) bool {
	return IsAuthError(err)
}

// 構造化エラー判定
func IsDomainError(err error) bool {
	var domainErr *DomainError
	return errors.As(err, &domainErr)
}
