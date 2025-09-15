package valueobject

import (
	"errors"
	"github.com/google/uuid"
)

// RoundID はラウンドIDを表すValue Object
type RoundID struct {
	value uuid.UUID
}

// NewRoundID は新しいRoundIDを作成する
func NewRoundID() RoundID {
	return RoundID{value: uuid.New()}
}

// NewRoundIDFromString は文字列からRoundIDを作成する
func NewRoundIDFromString(id string) (RoundID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return RoundID{}, errors.New("無効なラウンドID形式です")
	}
	return RoundID{value: parsed}, nil
}

// NewRoundIDFromUUID はuuid.UUIDからRoundIDを作成する
func NewRoundIDFromUUID(id uuid.UUID) RoundID {
	return RoundID{value: id}
}

// Value はuuid.UUIDの値を返す
func (r RoundID) Value() uuid.UUID {
	return r.value
}

// String は文字列表現を返す
func (r RoundID) String() string {
	return r.value.String()
}

// Equals は別のRoundIDと等価かを判定する
func (r RoundID) Equals(other RoundID) bool {
	return r.value == other.value
}

// IsEmpty は空のRoundIDかを判定する
func (r RoundID) IsEmpty() bool {
	return r.value == uuid.Nil
}