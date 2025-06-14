package errors

import (
	"errors"
	"fmt"
)

// =================================
// 技術的エラー統一
// =================================

var (
	// Database関連（PostgreSQL）
	ErrDatabaseConnection = errors.New("database connection failed")
	ErrDatabaseQuery      = errors.New("database query failed")
	ErrRecordNotFound     = errors.New("record not found")
	ErrUniqueConstraint   = errors.New("unique constraint violation")
	ErrTransactionFailed  = errors.New("transaction failed")
	ErrSQLExecution       = errors.New("SQL execution failed")

	// DynamoDB関連
	ErrDynamoDBConnection   = errors.New("DynamoDB connection failed")
	ErrDynamoDBOperation    = errors.New("DynamoDB operation failed")
	ErrDynamoDBItemNotFound = errors.New("DynamoDB item not found")
	ErrDynamoDBCondition    = errors.New("DynamoDB condition failed")

	// HTTP関連
	ErrHTTPRequest         = errors.New("HTTP request failed")
	ErrHTTPResponse        = errors.New("HTTP response invalid")
	ErrHTTPTimeout         = errors.New("HTTP timeout")
	ErrHTTPStatusCode      = errors.New("HTTP status code error")
	ErrHTTPRequestFailed   = errors.New("HTTP request failed")
	ErrHTTPResponseInvalid = errors.New("HTTP response invalid")

	// JWT/Auth Infrastructure関連
	ErrJWTParsingFailed    = errors.New("JWT parsing failed")
	ErrJWTTokenInvalid     = errors.New("JWT token invalid")
	ErrJWTTokenExpired     = errors.New("JWT token expired")
	ErrJWTSignatureInvalid = errors.New("JWT signature invalid")

	// JWKS関連
	ErrJWKSFetchFailed   = errors.New("JWKS fetch failed")
	ErrJWKSDecodeFailed  = errors.New("JWKS decode failed")
	ErrPublicKeyNotFound = errors.New("public key not found")
	ErrPublicKeyInvalid  = errors.New("public key invalid")
	ErrKeyIDNotFound     = errors.New("key ID not found")

	// SQS関連
	ErrSQSSendFailed       = errors.New("SQS send failed")
	ErrSQSConnectionFailed = errors.New("SQS connection failed")
	ErrSQSMessageInvalid   = errors.New("SQS message invalid")
	ErrSQSTimeout          = errors.New("SQS timeout")

	// Config関連
	ErrConfigMissing = errors.New("configuration missing")
	ErrConfigInvalid = errors.New("configuration invalid")

	// Cache関連
	ErrCacheExpired  = errors.New("cache expired")
	ErrCacheNotFound = errors.New("cache not found")
	ErrCacheInvalid  = errors.New("cache invalid")

	// File/IO関連
	ErrFileNotFound = errors.New("file not found")
	ErrFileRead     = errors.New("file read failed")
	ErrFileWrite    = errors.New("file write failed")
	ErrFileInvalid  = errors.New("file invalid")

	// Network関連
	ErrNetworkConnection  = errors.New("network connection failed")
	ErrNetworkTimeout     = errors.New("network timeout")
	ErrNetworkUnavailable = errors.New("network unavailable")

	// External Service関連
	ErrExternalServiceUnavailable = errors.New("external service unavailable")
	ErrExternalServiceTimeout     = errors.New("external service timeout")
	ErrExternalServiceError       = errors.New("external service error")
)

// =================================
// 構造化Infrastructureエラー
// =================================

type InfrastructureError struct {
	Component string // "database", "http", "sqs", "dynamodb", "jwt", "jwks"
	Operation string // "connect", "query", "send", "fetch", "parse"
	Message   string
	Cause     error
}

func (e *InfrastructureError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s.%s failed: %s (%v)", e.Component, e.Operation, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s.%s failed: %s", e.Component, e.Operation, e.Message)
}

func (e *InfrastructureError) Unwrap() error {
	return e.Cause
}

// =================================
// Infrastructure エラー作成ヘルパー関数
// =================================

func NewInfrastructureError(component, operation, message string, cause error) *InfrastructureError {
	return &InfrastructureError{
		Component: component,
		Operation: operation,
		Message:   message,
		Cause:     cause,
	}
}

// Database関連
func NewDatabaseError(operation string, cause error) error {
	return NewInfrastructureError("database", operation, "データベースエラー", cause)
}

func NewDatabaseConnectionError(cause error) error {
	return NewInfrastructureError("database", "connect", "データベース接続エラー", cause)
}

func NewDatabaseQueryError(cause error) error {
	return NewInfrastructureError("database", "query", "データベースクエリエラー", cause)
}

func NewTransactionError(cause error) error {
	return NewInfrastructureError("database", "transaction", "トランザクションエラー", cause)
}

func NewSQLExecutionError(cause error) error {
	return NewInfrastructureError("database", "execute", "SQL実行エラー", cause)
}

func NewUniqueConstraintError(cause error) error {
	return NewInfrastructureError("database", "constraint", "一意制約違反", cause)
}

// DynamoDB関連
func NewDynamoDBError(operation string, cause error) error {
	return NewInfrastructureError("dynamodb", operation, "DynamoDBエラー", cause)
}

func NewDynamoDBConnectionError(cause error) error {
	return NewInfrastructureError("dynamodb", "connect", "DynamoDB接続エラー", cause)
}

func NewDynamoDBOperationError(operation string, cause error) error {
	return NewInfrastructureError("dynamodb", operation, "DynamoDB操作エラー", cause)
}

func NewDynamoDBConditionError(cause error) error {
	return NewInfrastructureError("dynamodb", "condition", "DynamoDB条件エラー", cause)
}

// HTTP関連
func NewHTTPError(operation string, cause error) error {
	return NewInfrastructureError("http", operation, "HTTP通信エラー", cause)
}

func NewHTTPRequestError(cause error) error {
	return NewInfrastructureError("http", "request", "HTTPリクエストエラー", cause)
}

func NewHTTPResponseError(cause error) error {
	return NewInfrastructureError("http", "response", "HTTPレスポンスエラー", cause)
}

func NewHTTPTimeoutError(cause error) error {
	return NewInfrastructureError("http", "timeout", "HTTPタイムアウトエラー", cause)
}

func NewHTTPStatusCodeError(statusCode int, cause error) error {
	message := fmt.Sprintf("HTTPステータスコードエラー: %d", statusCode)
	return NewInfrastructureError("http", "status", message, cause)
}

// JWT関連
func NewJWTError(operation string, cause error) error {
	return NewInfrastructureError("jwt", operation, "JWT認証エラー", cause)
}

func NewJWTParsingError(cause error) error {
	return NewInfrastructureError("jwt", "parse", "JWT解析エラー", cause)
}

func NewJWTValidationError(cause error) error {
	return NewInfrastructureError("jwt", "validate", "JWT検証エラー", cause)
}

func NewJWTSignatureError(cause error) error {
	return NewInfrastructureError("jwt", "signature", "JWT署名エラー", cause)
}

// JWKS関連
func NewJWKSError(operation string, cause error) error {
	return NewInfrastructureError("jwks", operation, "JWKSエラー", cause)
}

func NewJWKSFetchError(cause error) error {
	return NewInfrastructureError("jwks", "fetch", "JWKS取得エラー", cause)
}

func NewJWKSDecodeError(cause error) error {
	return NewInfrastructureError("jwks", "decode", "JWKSデコードエラー", cause)
}

func NewPublicKeyError(cause error) error {
	return NewInfrastructureError("jwks", "public_key", "公開キーエラー", cause)
}

// SQS関連
func NewSQSError(operation string, cause error) error {
	return NewInfrastructureError("sqs", operation, "SQSエラー", cause)
}

func NewSQSSendError(cause error) error {
	return NewInfrastructureError("sqs", "send", "SQS送信エラー", cause)
}

func NewSQSConnectionError(cause error) error {
	return NewInfrastructureError("sqs", "connect", "SQS接続エラー", cause)
}

func NewSQSTimeoutError(cause error) error {
	return NewInfrastructureError("sqs", "timeout", "SQSタイムアウトエラー", cause)
}

// Config関連
func NewConfigError(key string, cause error) error {
	message := fmt.Sprintf("設定エラー: %s", key)
	return NewInfrastructureError("config", "load", message, cause)
}

func NewConfigMissingError(key string) error {
	message := fmt.Sprintf("必須設定が見つかりません: %s", key)
	return NewInfrastructureError("config", "missing", message, ErrConfigMissing)
}

func NewConfigInvalidError(key string, cause error) error {
	message := fmt.Sprintf("無効な設定値: %s", key)
	return NewInfrastructureError("config", "invalid", message, cause)
}

// Cache関連
func NewCacheError(operation string, cause error) error {
	return NewInfrastructureError("cache", operation, "キャッシュエラー", cause)
}

func NewCacheExpiredError() error {
	return NewInfrastructureError("cache", "expired", "キャッシュの有効期限切れ", ErrCacheExpired)
}

func NewCacheNotFoundError() error {
	return NewInfrastructureError("cache", "not_found", "キャッシュが見つかりません", ErrCacheNotFound)
}

// External Service関連
func NewExternalServiceError(service, operation string, cause error) error {
	message := fmt.Sprintf("外部サービスエラー: %s", service)
	return NewInfrastructureError("external", operation, message, cause)
}

func NewExternalServiceUnavailableError(service string) error {
	message := fmt.Sprintf("外部サービス利用不可: %s", service)
	return NewInfrastructureError("external", "unavailable", message, ErrExternalServiceUnavailable)
}

func NewExternalServiceTimeoutError(service string) error {
	message := fmt.Sprintf("外部サービスタイムアウト: %s", service)
	return NewInfrastructureError("external", "timeout", message, ErrExternalServiceTimeout)
}

// =================================
// エラー判定関数
// =================================

func IsInfrastructureError(err error) bool {
	var infraErr *InfrastructureError
	return errors.As(err, &infraErr) ||
		errors.Is(err, ErrDatabaseConnection) ||
		errors.Is(err, ErrHTTPRequest) ||
		errors.Is(err, ErrSQSSendFailed) ||
		errors.Is(err, ErrDynamoDBOperation) ||
		errors.Is(err, ErrJWTParsingFailed) ||
		errors.Is(err, ErrJWKSFetchFailed)
}

func IsDatabaseError(err error) bool {
	return errors.Is(err, ErrDatabaseConnection) ||
		errors.Is(err, ErrDatabaseQuery) ||
		errors.Is(err, ErrRecordNotFound) ||
		errors.Is(err, ErrTransactionFailed) ||
		errors.Is(err, ErrSQLExecution) ||
		errors.Is(err, ErrUniqueConstraint)
}

func IsDynamoDBError(err error) bool {
	return errors.Is(err, ErrDynamoDBConnection) ||
		errors.Is(err, ErrDynamoDBOperation) ||
		errors.Is(err, ErrDynamoDBItemNotFound) ||
		errors.Is(err, ErrDynamoDBCondition)
}

func IsHTTPError(err error) bool {
	return errors.Is(err, ErrHTTPRequest) ||
		errors.Is(err, ErrHTTPResponse) ||
		errors.Is(err, ErrHTTPTimeout) ||
		errors.Is(err, ErrHTTPStatusCode)
}

func IsJWTError(err error) bool {
	return errors.Is(err, ErrJWTParsingFailed) ||
		errors.Is(err, ErrJWTTokenInvalid) ||
		errors.Is(err, ErrJWTTokenExpired) ||
		errors.Is(err, ErrJWTSignatureInvalid)
}

func IsJWKSError(err error) bool {
	return errors.Is(err, ErrJWKSFetchFailed) ||
		errors.Is(err, ErrJWKSDecodeFailed) ||
		errors.Is(err, ErrPublicKeyNotFound) ||
		errors.Is(err, ErrPublicKeyInvalid) ||
		errors.Is(err, ErrKeyIDNotFound)
}

func IsSQSError(err error) bool {
	return errors.Is(err, ErrSQSSendFailed) ||
		errors.Is(err, ErrSQSConnectionFailed) ||
		errors.Is(err, ErrSQSTimeout)
}

func IsExternalServiceError(err error) bool {
	return errors.Is(err, ErrHTTPRequest) ||
		errors.Is(err, ErrSQSSendFailed) ||
		errors.Is(err, ErrDynamoDBOperation) ||
		errors.Is(err, ErrExternalServiceUnavailable) ||
		errors.Is(err, ErrExternalServiceTimeout)
}

func IsNetworkError(err error) bool {
	return errors.Is(err, ErrNetworkConnection) ||
		errors.Is(err, ErrNetworkTimeout) ||
		errors.Is(err, ErrNetworkUnavailable) ||
		errors.Is(err, ErrHTTPTimeout) ||
		errors.Is(err, ErrSQSTimeout)
}

func IsConfigError(err error) bool {
	return errors.Is(err, ErrConfigMissing) ||
		errors.Is(err, ErrConfigInvalid)
}

func IsCacheError(err error) bool {
	return errors.Is(err, ErrCacheExpired) ||
		errors.Is(err, ErrCacheNotFound) ||
		errors.Is(err, ErrCacheInvalid)
}

// 特定のエラータイプ判定
func IsTokenExpiredError(err error) bool {
	return errors.Is(err, ErrJWTTokenExpired)
}

func IsInvalidTokenError(err error) bool {
	return errors.Is(err, ErrJWTTokenInvalid) ||
		errors.Is(err, ErrJWTParsingFailed) ||
		errors.Is(err, ErrJWTSignatureInvalid)
}

func IsRecordNotFoundError(err error) bool {
	return errors.Is(err, ErrRecordNotFound) ||
		errors.Is(err, ErrDynamoDBItemNotFound)
}

func IsConnectionError(err error) bool {
	return errors.Is(err, ErrDatabaseConnection) ||
		errors.Is(err, ErrDynamoDBConnection) ||
		errors.Is(err, ErrSQSConnectionFailed) ||
		errors.Is(err, ErrNetworkConnection)
}

func IsTimeoutError(err error) bool {
	return errors.Is(err, ErrHTTPTimeout) ||
		errors.Is(err, ErrSQSTimeout) ||
		errors.Is(err, ErrNetworkTimeout) ||
		errors.Is(err, ErrExternalServiceTimeout)
}

// 構造化エラー判定
func IsInfrastructureErrorType(err error) bool {
	var infraErr *InfrastructureError
	return errors.As(err, &infraErr)
}
