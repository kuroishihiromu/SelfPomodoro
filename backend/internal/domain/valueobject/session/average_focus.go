package valueobject

import (
	"errors"
	"fmt"
)

// AverageFocus は平均集中度を表すValue Object
type AverageFocus struct {
	score float64
}

// NewAverageFocus は新しいAverageFocusを作成する
func NewAverageFocus(score float64) (AverageFocus, error) {
	if score < 0.0 || score > 100.0 {
		return AverageFocus{}, errors.New("平均集中度は0.0から100.0の間である必要があります")
	}
	return AverageFocus{score: score}, nil
}

// NewZeroAverageFocus は0のAverageFocusを作成する
func NewZeroAverageFocus() AverageFocus {
	return AverageFocus{score: 0.0}
}

// Score はスコア値を返す
func (a AverageFocus) Score() float64 {
	return a.score
}

// IsHigh は高い集中度かを判定する（80以上）
func (a AverageFocus) IsHigh() bool {
	return a.score >= 80.0
}

// IsMedium は中程度の集中度かを判定する（60-79.9）
func (a AverageFocus) IsMedium() bool {
	return a.score >= 60.0 && a.score < 80.0
}

// IsLow は低い集中度かを判定する（60未満）
func (a AverageFocus) IsLow() bool {
	return a.score < 60.0
}

// ToQualityLevel は品質レベルを返す
func (a AverageFocus) ToQualityLevel() string {
	if a.IsHigh() {
		return "高"
	} else if a.IsMedium() {
		return "中"
	} else {
		return "低"
	}
}

// IsProductiveSession は生産的なセッションかを判定する（70以上）
func (a AverageFocus) IsProductiveSession() bool {
	return a.score >= 70.0
}

// CalculateWeightedAverage は重み付き平均を計算する（セッション統計用）
func CalculateWeightedAverageFocus(scores []float64, weights []int) (AverageFocus, error) {
	if len(scores) != len(weights) {
		return AverageFocus{}, errors.New("スコアと重みの配列長が一致しません")
	}
	if len(scores) == 0 {
		return NewZeroAverageFocus(), nil
	}

	var totalWeightedScore float64
	var totalWeight int
	for i, score := range scores {
		totalWeightedScore += score * float64(weights[i])
		totalWeight += weights[i]
	}

	if totalWeight == 0 {
		return NewZeroAverageFocus(), nil
	}

	average := totalWeightedScore / float64(totalWeight)
	return NewAverageFocus(average)
}

// String は文字列表現を返す
func (a AverageFocus) String() string {
	return fmt.Sprintf("%.1f%%", a.score)
}

// Equals は別のAverageFocusと等価かを判定する
func (a AverageFocus) Equals(other AverageFocus) bool {
	// 浮動小数点の比較は0.1の精度で行う
	return fmt.Sprintf("%.1f", a.score) == fmt.Sprintf("%.1f", other.score)
}