package statistics

import (
	"errors"
	"fmt"
)

// AverageFocusScore は統計の平均集中度スコアを表すValue Object
type AverageFocusScore struct {
	score float64
}

// NewAverageFocusScore は新しいAverageFocusScoreを作成する
func NewAverageFocusScore(score float64) (AverageFocusScore, error) {
	if score < 0.0 || score > 100.0 {
		return AverageFocusScore{}, errors.New("平均集中度スコアは0.0から100.0の間である必要があります")
	}
	return AverageFocusScore{score: score}, nil
}

// NewZeroAverageFocusScore は0のAverageFocusScoreを作成する
func NewZeroAverageFocusScore() AverageFocusScore {
	return AverageFocusScore{score: 0.0}
}

// Score はスコア値を返す
func (a AverageFocusScore) Score() float64 {
	return a.score
}

// IsZero は0スコアかを判定する
func (a AverageFocusScore) IsZero() bool {
	return a.score == 0.0
}

// IsExcellent は優秀なスコアかを判定する（90以上）
func (a AverageFocusScore) IsExcellent() bool {
	return a.score >= 90.0
}

// IsHigh は高いスコアかを判定する（80以上）
func (a AverageFocusScore) IsHigh() bool {
	return a.score >= 80.0
}

// IsMedium は中程度のスコアかを判定する（60-79.9）
func (a AverageFocusScore) IsMedium() bool {
	return a.score >= 60.0 && a.score < 80.0
}

// IsLow は低いスコアかを判定する（40-59.9）
func (a AverageFocusScore) IsLow() bool {
	return a.score >= 40.0 && a.score < 60.0
}

// IsPoor は悪いスコアかを判定する（40未満）
func (a AverageFocusScore) IsPoor() bool {
	return a.score < 40.0
}

// ToQualityLevel は品質レベルを返す
func (a AverageFocusScore) ToQualityLevel() string {
	if a.IsExcellent() {
		return "優秀"
	} else if a.IsHigh() {
		return "高"
	} else if a.IsMedium() {
		return "中"
	} else if a.IsLow() {
		return "低"
	} else {
		return "要改善"
	}
}

// ToQualityGrade はA-Fの成績評価を返す
func (a AverageFocusScore) ToQualityGrade() string {
	if a.score >= 90 {
		return "A"
	} else if a.score >= 80 {
		return "B"
	} else if a.score >= 70 {
		return "C"
	} else if a.score >= 60 {
		return "D"
	} else {
		return "F"
	}
}

// IsProductiveDay は生産的な日かを判定する（70以上）
func (a AverageFocusScore) IsProductiveDay() bool {
	return a.score >= 70.0
}

// IsProductiveWeek は生産的な週かを判定する（75以上）
func (a AverageFocusScore) IsProductiveWeek() bool {
	return a.score >= 75.0
}

// CalculateWeightedAverage は重み付き平均を計算する
func CalculateWeightedAverageFocusScore(scores []float64, weights []int) (AverageFocusScore, error) {
	if len(scores) != len(weights) {
		return AverageFocusScore{}, errors.New("スコアと重みの配列長が一致しません")
	}
	if len(scores) == 0 {
		return NewZeroAverageFocusScore(), nil
	}

	var totalWeightedScore float64
	var totalWeight int
	for i, score := range scores {
		if score < 0 || score > 100 {
			return AverageFocusScore{}, fmt.Errorf("無効なスコア値: %f", score)
		}
		totalWeightedScore += score * float64(weights[i])
		totalWeight += weights[i]
	}

	if totalWeight == 0 {
		return NewZeroAverageFocusScore(), nil
	}

	average := totalWeightedScore / float64(totalWeight)
	return NewAverageFocusScore(average)
}

// UpdateAverage は新しいスコアで平均を更新する
func (a AverageFocusScore) UpdateAverage(newScore float64, currentCount int) (AverageFocusScore, error) {
	if newScore < 0 || newScore > 100 {
		return AverageFocusScore{}, fmt.Errorf("無効なスコア値: %f", newScore)
	}
	if currentCount < 0 {
		return AverageFocusScore{}, errors.New("現在のカウントは0以上である必要があります")
	}

	totalScore := a.score * float64(currentCount)
	totalScore += newScore
	newCount := currentCount + 1
	newAverage := totalScore / float64(newCount)

	return NewAverageFocusScore(newAverage)
}

// String は文字列表現を返す
func (a AverageFocusScore) String() string {
	if a.IsZero() {
		return "未記録"
	}
	return fmt.Sprintf("%.1f%%", a.score)
}

// Equals は別のAverageFocusScoreと等価かを判定する
func (a AverageFocusScore) Equals(other AverageFocusScore) bool {
	// 浮動小数点の比較は0.1の精度で行う
	return fmt.Sprintf("%.1f", a.score) == fmt.Sprintf("%.1f", other.score)
}

// Compare は他のAverageFocusScoreとの大小を比較する
// 戻り値: -1 (小さい), 0 (等しい), 1 (大きい)
func (a AverageFocusScore) Compare(other AverageFocusScore) int {
	if a.score < other.score {
		return -1
	}
	if a.score > other.score {
		return 1
	}
	return 0
}