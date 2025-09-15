package optimization

import (
	"errors"
	"fmt"
	"time"

	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
	roundVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/round"
)

// RoundOptimizationRecord はラウンド最適化記録を表すValue Object
type RoundOptimizationRecord struct {
	userID     userVO.UserID
	timestamp  time.Time
	workTime   roundVO.WorkTime
	breakTime  roundVO.BreakTime
	focusScore roundVO.FocusScore
}

// NewRoundOptimizationRecord は新しいラウンド最適化記録を作成する
func NewRoundOptimizationRecord(userID userVO.UserID, workTime roundVO.WorkTime, breakTime roundVO.BreakTime, focusScore roundVO.FocusScore) (RoundOptimizationRecord, error) {
	// 作業時間と休憩時間の論理チェック
	if workTime.Minutes() <= 0 {
		return RoundOptimizationRecord{}, errors.New("作業時間は1分以上である必要があります")
	}
	
	if breakTime.Minutes() < 0 {
		return RoundOptimizationRecord{}, errors.New("休憩時間は0分以上である必要があります")
	}

	return RoundOptimizationRecord{
		userID:     userID,
		timestamp:  time.Now(),
		workTime:   workTime,
		breakTime:  breakTime,
		focusScore: focusScore,
	}, nil
}

// NewRoundOptimizationRecordWithTime は指定時刻でラウンド最適化記録を作成する（テスト・サンプルデータ用）
func NewRoundOptimizationRecordWithTime(userID userVO.UserID, timestamp time.Time, workTime roundVO.WorkTime, breakTime roundVO.BreakTime, focusScore roundVO.FocusScore) (RoundOptimizationRecord, error) {
	if workTime.Minutes() <= 0 {
		return RoundOptimizationRecord{}, errors.New("作業時間は1分以上である必要があります")
	}
	
	if breakTime.Minutes() < 0 {
		return RoundOptimizationRecord{}, errors.New("休憩時間は0分以上である必要があります")
	}

	return RoundOptimizationRecord{
		userID:     userID,
		timestamp:  timestamp,
		workTime:   workTime,
		breakTime:  breakTime,
		focusScore: focusScore,
	}, nil
}

// UserID はユーザーIDを返す
func (r RoundOptimizationRecord) UserID() userVO.UserID {
	return r.userID
}

// Timestamp はタイムスタンプを返す
func (r RoundOptimizationRecord) Timestamp() time.Time {
	return r.timestamp
}

// WorkTime は作業時間を返す
func (r RoundOptimizationRecord) WorkTime() roundVO.WorkTime {
	return r.workTime
}

// BreakTime は休憩時間を返す
func (r RoundOptimizationRecord) BreakTime() roundVO.BreakTime {
	return r.breakTime
}

// FocusScore は集中度スコアを返す
func (r RoundOptimizationRecord) FocusScore() roundVO.FocusScore {
	return r.focusScore
}

// Equals は別のRoundOptimizationRecordと等価かを判定する
func (r RoundOptimizationRecord) Equals(other RoundOptimizationRecord) bool {
	return r.userID.Equals(other.userID) &&
		r.timestamp.Equal(other.timestamp) &&
		r.workTime.Equals(other.workTime) &&
		r.breakTime.Equals(other.breakTime) &&
		r.focusScore.Equals(other.focusScore)
}

// IsHighPerformance は高いパフォーマンスの最適化記録かを判定する
func (r RoundOptimizationRecord) IsHighPerformance() bool {
	return r.focusScore.IsHigh()
}

// IsEffectiveOptimization は効果的な最適化かを判定する
func (r RoundOptimizationRecord) IsEffectiveOptimization() bool {
	// 集中度が中程度以上で、時間設定が合理的
	return r.focusScore.Value() >= 60 && 
		   r.workTime.Minutes() >= 15 && r.workTime.Minutes() <= 60 &&
		   r.breakTime.Minutes() <= 15
}

// GetOptimizationQuality は最適化の品質レベルを返す
func (r RoundOptimizationRecord) GetOptimizationQuality() string {
	if r.focusScore.IsHigh() {
		return "高品質"
	} else if r.focusScore.IsMedium() {
		return "中品質"
	} else {
		return "低品質"
	}
}

// ToTimestamp は永続化用のタイムスタンプ文字列を返す
func (r RoundOptimizationRecord) ToTimestamp() string {
	return r.timestamp.Format(time.RFC3339)
}

// String は文字列表現を返す
func (r RoundOptimizationRecord) String() string {
	return fmt.Sprintf("RoundOptimization[User=%s, Time=%s, Work=%s, Break=%s, Focus=%s]",
		r.userID.String()[:8], 
		r.timestamp.Format("2006-01-02 15:04"),
		r.workTime.String(),
		r.breakTime.String(),
		r.focusScore.String())
}