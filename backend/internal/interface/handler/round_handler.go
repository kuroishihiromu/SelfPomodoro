package handler

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/postgres"
	"github.com/tsunakit99/selfpomodoro/internal/usecase"
)

// RoundHandler はラウンド関連のエンドポイントを処理するハンドラー
type RoundHandler struct {
	roundUseCase usecase.RoundUseCase
	logger       logger.Logger
	validator    *validator.Validate
}

// NewRoundHandler は新しいRoundHandlerインスタンスを作成する
func NewRoundHandler(roundUseCase usecase.RoundUseCase, logger logger.Logger) *RoundHandler {
	return &RoundHandler{
		roundUseCase: roundUseCase,
		logger:       logger,
		validator:    validator.New(),
	}
}

// getUserIDFromContext はコンテキストからユーザーIDを取得する
func (h *RoundHandler) getUserIDFromContext(c echo.Context) (uuid.UUID, error) {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return uuid.Nil, errors.New("ユーザーIDが見つかりません")
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, errors.New("無効なユーザーID")
	}
	return id, nil
}

// StartRound はラウンド開始エンドポイントを処理する
// POST /sessions/:session_id/rounds
func (h *RoundHandler) StartRound(c echo.Context) error {
	// ユーザーIDの取得
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// セッションIDのパース
	sessionID, err := uuid.Parse(c.Param("session_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "無効なセッションID"})
	}

	// リクエストボディのバインド（一旦空でもいい設計）
	var req model.RoundCreateRequest
	if err := c.Bind(&req); err != nil {
		// TODO: RoundCreateRequestの有無を決める
	}

	// ユースケースを呼び出してラウンドを開始
	roundResponse, err := h.roundUseCase.StartRound(c.Request().Context(), sessionID, userID, &req)
	if err != nil {
		if errors.Is(err, postgres.ErrSessionNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "セッションが見つかりません"})
		}
		if errors.Is(err, postgres.ErrSessionAccessDenied) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "このセッションへのアクセス権限がありません"})
		}

		h.logger.Errorf("ラウンド開始エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ラウンド開始に失敗しました"})
	}
	return c.JSON(http.StatusCreated, roundResponse)
}

// GetRound はラウンド取得エンドポイントを処理する
// GET /rounds/:round_id
func (h *RoundHandler) GetRound(c echo.Context) error {
	// ラウンドIDのパース
	roundID, err := uuid.Parse(c.Param("round_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "無効なラウンドID"})
	}

	// ユースケースを呼び出してラウンドを取得
	roundResponse, err := h.roundUseCase.GetRound(c.Request().Context(), roundID)
	if err != nil {
		if errors.Is(err, postgres.ErrRoundNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "ラウンドが見つかりません"})
		}

		h.logger.Errorf("ラウンド取得エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ラウンド取得に失敗しました"})
	}

	return c.JSON(http.StatusOK, roundResponse)
}

// GetAllRoundsBySessionID はセッションのラウンド一覧取得エンドポイントを処理する
// GET /sessions/:session_id/rounds
func (h *RoundHandler) GetAllRoundsBySessionID(c echo.Context) error {
	// セッションIDのパース
	sessionID, err := uuid.Parse(c.Param("session_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "無効なセッションID"})
	}

	// ユースケースを呼び出してラウンド一覧を取得
	roundsResponse, err := h.roundUseCase.GetAllRoundsBySessionID(c.Request().Context(), sessionID)
	if err != nil {
		h.logger.Errorf("ラウンド一覧取得エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ラウンド一覧の取得に失敗しました"})
	}
	return c.JSON(http.StatusOK, roundsResponse)
}

// CompleteRound はラウンド完了エンドポイントを処理する
// POST /rounds/:round_id/complete
func (h *RoundHandler) CompleteRound(c echo.Context) error {
	// ユーザーIDの取得
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// ラウンドIDのパース
	roundID, err := uuid.Parse(c.Param("round_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "無効なラウンドID"})
	}

	// リクエストボディのバインド
	var req model.RoundCompleteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "無効なリクエスト形式"})
	}

	// バリデーション
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "集中度スコアは0から100の間である必要があります"})
	}

	// ユースケースを呼び出してラウンドを完了
	roundResponse, err := h.roundUseCase.CompleteRound(c.Request().Context(), roundID, userID, &req)
	if err != nil {
		if errors.Is(err, postgres.ErrRoundNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "ラウンドが見つかりません"})
		}

		h.logger.Errorf("ラウンド完了エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ラウンド完了に失敗しました"})
	}

	return c.JSON(http.StatusOK, roundResponse)
}

// AbortRound はラウンド中止エンドポイントを処理する
// POST /rounds/:round_id/abort
func (h *RoundHandler) AbortRound(c echo.Context) error {
	// ユーザーIDの取得
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// ラウンドIDのパース
	roundID, err := uuid.Parse(c.Param("round_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "無効なラウンドID"})
	}

	// ユースケースを呼び出してラウンドを中止
	roundResponse, err := h.roundUseCase.AbortRound(c.Request().Context(), roundID, userID)
	if err != nil {
		if errors.Is(err, postgres.ErrRoundNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "ラウンドが見つかりません"})
		}

		h.logger.Errorf("ラウンド中止エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ラウンド中止に失敗しました"})
	}

	return c.JSON(http.StatusOK, roundResponse)
}
