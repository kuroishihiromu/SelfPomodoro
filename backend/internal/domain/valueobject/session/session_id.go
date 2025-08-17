package valueobject

import (
	"errors"
	"github.com/google/uuid"
)

// SessionID はセッションIDを表すValue Object
type SessionID struct {
	value uuid.UUID
}

// NewSessionID は新しいSessionIDを作成する
func NewSessionID() SessionID {
	return SessionID{value: uuid.New()}
}

// NewSessionIDFromString は文字列からSessionIDを作成する
func NewSessionIDFromString(id string) (SessionID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return SessionID{}, errors.New("無効なセッションID形式です")
	}
	return SessionID{value: parsed}, nil
}

// NewSessionIDFromUUID はuuid.UUIDからSessionIDを作成する
func NewSessionIDFromUUID(id uuid.UUID) SessionID {
	return SessionID{value: id}
}

// Value はuuid.UUIDの値を返す
func (s SessionID) Value() uuid.UUID {
	return s.value
}

// String は文字列表現を返す
func (s SessionID) String() string {
	return s.value.String()
}

// Equals は別のSessionIDと等価かを判定する
func (s SessionID) Equals(other SessionID) bool {
	return s.value == other.value
}

// IsEmpty は空のSessionIDかを判定する
func (s SessionID) IsEmpty() bool {
	return s.value == uuid.Nil
}