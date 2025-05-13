package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// 開発環境用の簡易認証ミドルウェア
// TODO: 本番環境ではCognitoを使用する
func AuthMiddleware() echo.MiddlewareFunc {
	// ハンドラーをラップする関数
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// 実際に処理されるラップされたハンドラ本体
		return func(c echo.Context) error {
			// Authorizationヘッダーの取得
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "認証情報が必要です"})
			}

			// Bearerトークンのチェック
			parts := strings.Split(auth, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "無効な認証情報"})
			}

			token := parts[1]

			// 開発環境用の簡易トークン検証
			// TODO: 本番環境ではCognitoを使用する
			if token == "dev-token" {
				// 開発用固定ユーザーID
				c.Set("user_id", "00000000-0000-0000-0000-000000000001")
				// 通過したら次のハンドラーに進む
				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "無効なトークンです"})
		}
	}
}
