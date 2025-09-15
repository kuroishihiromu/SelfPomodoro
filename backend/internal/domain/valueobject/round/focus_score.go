package valueobject

import (
	"fmt"
)

// 最適化関連定数
const (
	MaxFocusScore = 100
	MinFocusScore = 0
)

// FocusScore は集中度スコアを表すValue Object
type FocusScore struct {
	value int
}

// NewFocusScore は新しいFocusScoreを作成する
func NewFocusScore(score int) (FocusScore, error) {
	if score < MinFocusScore || score > MaxFocusScore {
		return FocusScore{}, fmt.Errorf("集中度スコアは%dから%dの間で設定してください", MinFocusScore, MaxFocusScore)
	}
	return FocusScore{value: score}, nil
}

// MustNewFocusScore はパニックを起こす可能性があるFocusScore作成（テスト用）
func MustNewFocusScore(score int) FocusScore {
	fs, err := NewFocusScore(score)
	if err != nil {
		panic(err)
	}
	return fs
}

// Value はintの値を返す
func (f FocusScore) Value() int {
	return f.value
}

// Equals は別のFocusScoreと等価かを判定する
func (f FocusScore) Equals(other FocusScore) bool {
	return f.value == other.value
}

// IsHigh は高い集中度かを判定する（80以上）
func (f FocusScore) IsHigh() bool {
	return f.value >= 80
}

// IsMedium は中程度の集中度かを判定する（60-79）
func (f FocusScore) IsMedium() bool {
	return f.value >= 60 && f.value < 80
}

// IsLow は低い集中度かを判定する（59以下）
func (f FocusScore) IsLow() bool {
	return f.value < 60
}

// ToQualityString は品質の文字列表現を返す
func (f FocusScore) ToQualityString() string {
	if f.IsHigh() {
		return "高"
	} else if f.IsMedium() {
		return "中"
	} else {
		return "低"
	}
}

// String は文字列表現を返す
func (f FocusScore) String() string {
	return fmt.Sprintf("%d", f.value)
}