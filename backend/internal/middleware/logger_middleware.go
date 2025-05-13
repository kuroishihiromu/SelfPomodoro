package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// LogggerMiddleware はリクエストとレスポンスをログに記録するミドルウェア
func LogggerMiddleware(logger logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()

			// リクエストのログを記録
			logger.Infof("REQUEST: %s %s %s", req.Method, req.URL.Path, req.RemoteAddr)

			// 次のハンドラーを呼び出す
			err := next(c)

			// レスポンスのログを記録
			stop := time.Now()
			latency := stop.Sub(start)
			logger.Infof("RESPONSE: %s %s %d %s", req.Method, req.URL.Path, res.Status, latency)

			return err
		}
	}
}
