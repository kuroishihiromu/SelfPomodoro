package statistics

import (
	"errors"
	"fmt"
)

// SessionCount は統計のセッション数を表すValue Object
type SessionCount struct {
	count int
}

// NewSessionCount は新しいSessionCountを作成する
func NewSessionCount(count int) (SessionCount, error) {
	if count < 0 {
		return SessionCount{}, errors.New("セッション数は0以上である必要があります")
	}
	// 制限なし（ユーザーの要求通り）
	return SessionCount{count: count}, nil
}

// NewZeroSessionCount は0のSessionCountを作成する
func NewZeroSessionCount() SessionCount {
	return SessionCount{count: 0}
}

// Count は数値を返す
func (s SessionCount) Count() int {
	return s.count
}

// Increment はセッション数を1増やす
func (s SessionCount) Increment() (SessionCount, error) {
	return NewSessionCount(s.count + 1)
}

// Add は指定した数だけセッション数を増やす
func (s SessionCount) Add(sessions int) (SessionCount, error) {
	return NewSessionCount(s.count + sessions)
}

// IsZero は0セッションかを判定する
func (s SessionCount) IsZero() bool {
	return s.count == 0
}

// IsMinimal は最小セッション数（1セッション）かを判定する
func (s SessionCount) IsMinimal() bool {
	return s.count == 1
}

// IsActiveDay はアクティブな日（2セッション以上）かを判定する
func (s SessionCount) IsActiveDay() bool {
	return s.count >= 2
}

// IsProductiveDay は生産的な日（5セッション以上）かを判定する
func (s SessionCount) IsProductiveDay() bool {
	return s.count >= 5
}

// IsIntensiveDay は集中的な日（10セッション以上）かを判定する
func (s SessionCount) IsIntensiveDay() bool {
	return s.count >= 10
}

// ToActivityLevel はセッション数からアクティビティレベルを返す
func (s SessionCount) ToActivityLevel() string {
	if s.count == 0 {
		return "非アクティブ"
	} else if s.count == 1 {
		return "軽微"
	} else if s.count <= 3 {
		return "普通"
	} else if s.count <= 7 {
		return "アクティブ"
	} else {
		return "超アクティブ"
	}
}

// String は文字列表現を返す
func (s SessionCount) String() string {
	if s.count == 0 {
		return "セッションなし"
	}
	return fmt.Sprintf("%dセッション", s.count)
}

// Equals は別のSessionCountと等価かを判定する
func (s SessionCount) Equals(other SessionCount) bool {
	return s.count == other.count
}

// Compare は他のSessionCountとの大小を比較する
// 戻り値: -1 (小さい), 0 (等しい), 1 (大きい)
func (s SessionCount) Compare(other SessionCount) int {
	if s.count < other.count {
		return -1
	}
	if s.count > other.count {
		return 1
	}
	return 0
}