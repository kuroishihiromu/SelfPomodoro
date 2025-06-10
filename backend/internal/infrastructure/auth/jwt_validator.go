package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	authErrors "github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/auth/errors"
)

// JWKS は JSON Web Key Set を表す構造体
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// JWK は JSON Web Key を表す構造体
type JWK struct {
	Kty string `json:"kty"` // Key Type
	Use string `json:"use"` // Public Key Use
	Kid string `json:"kid"` // Key ID
	N   string `json:"n"`   // Modulus
	E   string `json:"e"`   // Exponent
	Alg string `json:"alg"` // Algorithm
}

// CognitoJWTValidator はCognito JWTを検証するためのバリデーター
type CognitoJWTValidator struct {
	userPoolID string
	clientID   string
	region     string
	issuerURL  string
	jwksURL    string

	// 公開キーキャッシュ
	keyCache     map[string]*rsa.PublicKey
	cacheExpiry  time.Time
	cacheMutex   sync.RWMutex
	cacheTimeout time.Duration

	// HTTPクライアント
	httpClient *http.Client
	logger     logger.Logger
}

// CognitoJWTValidatorConfig はバリデーターの設定
type CognitoJWTValidatorConfig struct {
	UserPoolID   string
	ClientID     string
	Region       string
	CacheTimeout time.Duration
	HTTPTimeout  time.Duration
}

// NewCognitoJWTValidator は新しいCognito JWT バリデーターを作成する
func NewCognitoJWTValidator(config *CognitoJWTValidatorConfig, logger logger.Logger) *CognitoJWTValidator {
	// デフォルト設定
	if config.CacheTimeout == 0 {
		config.CacheTimeout = 1 * time.Hour // 1時間キャッシュ
	}
	if config.HTTPTimeout == 0 {
		config.HTTPTimeout = 10 * time.Second
	}

	// URLs構築
	issuerURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", config.Region, config.UserPoolID)
	jwksURL := fmt.Sprintf("%s/.well-known/jwks.json", issuerURL)

	return &CognitoJWTValidator{
		userPoolID:   config.UserPoolID,
		clientID:     config.ClientID,
		region:       config.Region,
		issuerURL:    issuerURL,
		jwksURL:      jwksURL,
		keyCache:     make(map[string]*rsa.PublicKey),
		cacheTimeout: config.CacheTimeout,
		httpClient: &http.Client{
			Timeout: config.HTTPTimeout,
		},
		logger: logger,
	}
}

// ValidateJWT はJWTトークンを検証し、クレームを返す（インフラエラー版）
func (v *CognitoJWTValidator) ValidateJWT(tokenString string) (*CognitoClaims, error) {
	v.logger.Infof("JWT検証開始: トークン長=%d", len(tokenString))

	// トークン形式の事前チェック
	if tokenString == "" {
		return nil, authErrors.ErrJWTTokenNotFound
	}

	// Bearer プレフィックスを除去
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	// JWTをパースしてヘッダーを取得
	token, err := jwt.ParseWithClaims(tokenString, &CognitoClaims{}, v.keyFunc)
	if err != nil {
		v.logger.Errorf("JWT解析エラー: %v", err)
		if strings.Contains(err.Error(), "token is expired") {
			return nil, authErrors.ErrJWTTokenExpired
		}
		return nil, authErrors.NewJWTParsingError(err)
	}

	// クレームの取得と検証
	claims, ok := token.Claims.(*CognitoClaims)
	if !ok {
		return nil, authErrors.ErrJWTTokenInvalid
	}

	// 発行者の検証
	if claims.Issuer != v.issuerURL {
		v.logger.Errorf("無効な発行者: expected=%s, actual=%s", v.issuerURL, claims.Issuer)
		return nil, authErrors.ErrJWTInvalidIssuer
	}

	// オーディエンスの検証（IDトークンの場合）
	if claims.IsIDToken() && claims.Audience != v.clientID {
		v.logger.Errorf("無効なオーディエンス: expected=%s, actual=%s", v.clientID, claims.Audience)
		return nil, authErrors.ErrJWTInvalidAudience
	}

	// クレームの有効性検証
	if err := claims.IsValid(); err != nil {
		v.logger.Errorf("クレーム検証エラー: %v", err)
		// Infrastructureエラーをそのまま返す（UseCase層で変換）
		return nil, err
	}

	v.logger.Infof("JWT検証成功: %s", claims.ToLogString())
	return claims, nil
}

// keyFunc はJWTパーサー用のキー取得関数
func (v *CognitoJWTValidator) keyFunc(token *jwt.Token) (interface{}, error) {
	// アルゴリズムの確認
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, authErrors.NewJWTParsingError(fmt.Errorf("予期しない署名方式: %v", token.Header["alg"]))
	}

	// Key IDの取得
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, authErrors.ErrPublicKeyNotFound
	}

	// 公開キーの取得
	publicKey, err := v.getPublicKey(kid)
	if err != nil {
		return nil, err // インフラエラーをそのまま返す
	}

	return publicKey, nil
}

// getPublicKey は指定されたKey IDの公開キーを取得する（キャッシュ機能付き）
func (v *CognitoJWTValidator) getPublicKey(kid string) (*rsa.PublicKey, error) {
	// キャッシュから取得を試行
	v.cacheMutex.RLock()
	if publicKey, exists := v.keyCache[kid]; exists && time.Now().Before(v.cacheExpiry) {
		v.cacheMutex.RUnlock()
		v.logger.Debugf("公開キーをキャッシュから取得: kid=%s", kid)
		return publicKey, nil
	}
	v.cacheMutex.RUnlock()

	// キャッシュにない場合は取得
	v.logger.Infof("JWKS取得開始: %s", v.jwksURL)
	return v.fetchAndCachePublicKey(kid)
}

// fetchAndCachePublicKey はJWKSエンドポイントから公開キーを取得してキャッシュする
func (v *CognitoJWTValidator) fetchAndCachePublicKey(kid string) (*rsa.PublicKey, error) {
	// JWKS取得
	resp, err := v.httpClient.Get(v.jwksURL)
	if err != nil {
		return nil, authErrors.NewHTTPError(fmt.Errorf("JWKS取得リクエストエラー: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, authErrors.NewHTTPError(fmt.Errorf("JWKS取得HTTPエラー: %d", resp.StatusCode))
	}

	// JWKSデコード
	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, authErrors.NewJWKSError(fmt.Errorf("JWKSデコードエラー: %w", err))
	}

	v.logger.Infof("JWKS取得成功: %d keys", len(jwks.Keys))

	// 全キーをキャッシュに保存
	v.cacheMutex.Lock()
	defer v.cacheMutex.Unlock()

	v.keyCache = make(map[string]*rsa.PublicKey)
	v.cacheExpiry = time.Now().Add(v.cacheTimeout)

	var targetKey *rsa.PublicKey
	for _, jwk := range jwks.Keys {
		publicKey, err := v.jwkToRSAPublicKey(&jwk)
		if err != nil {
			v.logger.Warnf("JWK変換エラー (kid=%s): %v", jwk.Kid, err)
			continue
		}

		v.keyCache[jwk.Kid] = publicKey
		if jwk.Kid == kid {
			targetKey = publicKey
		}
	}

	if targetKey == nil {
		return nil, authErrors.NewPublicKeyError(fmt.Errorf("指定されたKey ID(%s)が見つかりません", kid))
	}

	v.logger.Infof("公開キーキャッシュ更新完了: %d keys cached", len(v.keyCache))
	return targetKey, nil
}

// jwkToRSAPublicKey はJWKをRSA公開キーに変換する
func (v *CognitoJWTValidator) jwkToRSAPublicKey(jwk *JWK) (*rsa.PublicKey, error) {
	// Base64URLデコード
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("modulus デコードエラー: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("exponent デコードエラー: %w", err)
	}

	// RSA公開キー構築
	var e int
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	n := new(big.Int).SetBytes(nBytes)

	return &rsa.PublicKey{
		N: n,
		E: e,
	}, nil
}

// HealthCheck はバリデーターの接続確認を行う
func (v *CognitoJWTValidator) HealthCheck(ctx context.Context) error {
	// JWKS エンドポイントへの接続確認
	req, err := http.NewRequestWithContext(ctx, "GET", v.jwksURL, nil)
	if err != nil {
		return authErrors.NewHTTPError(fmt.Errorf("JWKS HealthCheck リクエスト作成エラー: %w", err))
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return authErrors.NewHTTPError(fmt.Errorf("JWKS HealthCheck 接続エラー: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return authErrors.NewHTTPError(fmt.Errorf("JWKS HealthCheck HTTPエラー: %d", resp.StatusCode))
	}

	v.logger.Info("CognitoJWTValidator HealthCheck 成功")
	return nil
}

// ClearCache はキャッシュをクリアする（テスト用）
func (v *CognitoJWTValidator) ClearCache() {
	v.cacheMutex.Lock()
	defer v.cacheMutex.Unlock()

	v.keyCache = make(map[string]*rsa.PublicKey)
	v.cacheExpiry = time.Time{}
	v.logger.Debug("公開キーキャッシュクリア完了")
}
