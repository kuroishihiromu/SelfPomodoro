package statistics

import (
	"fmt"

	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
)

// FocusScore は集中度スコアを表すValue Object
type FocusScore struct {
	value float64
}

// NewFocusScore は新しい集中度スコアを作成する
func NewFocusScore(value float64) (FocusScore, error) {
	if value < 0 || value > 100 {
		return FocusScore{}, appErrors.NewValidationError(fmt.Sprintf("集中度スコアは0-100の範囲である必要があります: %.2f", value))
	}
	return FocusScore{value: value}, nil
}

// Value は集中度スコアの値を返す
func (fs FocusScore) Value() float64 {
	return fs.value
}

// IsHigh は高い集中度かを判定する（80以上）
func (fs FocusScore) IsHigh() bool {
	return fs.value >= 80
}

// IsMedium は中程度の集中度かを判定する（60以上80未満）
func (fs FocusScore) IsMedium() bool {
	return fs.value >= 60 && fs.value < 80
}

// IsLow は低い集中度かを判定する（60未満）
func (fs FocusScore) IsLow() bool {
	return fs.value < 60
}

// GetLevel は集中度レベルを文字列で返す
func (fs FocusScore) GetLevel() string {
	if fs.IsHigh() {
		return "高"
	} else if fs.IsMedium() {
		return "中"
	} else {
		return "低"
	}
}

// Equals は別のFocusScoreと等価かを判定する
func (fs FocusScore) Equals(other FocusScore) bool {
	return fs.value == other.value
}

// String は文字列表現を返す
func (fs FocusScore) String() string {
	return fmt.Sprintf("%.1f", fs.value)
}