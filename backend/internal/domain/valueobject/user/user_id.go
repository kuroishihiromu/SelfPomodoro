package user

import (
	"errors"
	"github.com/google/uuid"
)

// UserID はユーザーIDを表すValue Object
type UserID struct {
	value uuid.UUID
}

// NewUserID は新しいUserIDを作成する
func NewUserID() UserID {
	return UserID{value: uuid.New()}
}

// NewUserIDFromString は文字列からUserIDを作成する
func NewUserIDFromString(id string) (UserID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return UserID{}, errors.New("無効なユーザーID形式です")
	}
	return UserID{value: parsed}, nil
}

// NewUserIDFromUUID はuuid.UUIDからUserIDを作成する
func NewUserIDFromUUID(id uuid.UUID) UserID {
	return UserID{value: id}
}

// Value はuuid.UUIDの値を返す
func (u UserID) Value() uuid.UUID {
	return u.value
}

// String は文字列表現を返す
func (u UserID) String() string {
	return u.value.String()
}

// Equals は別のUserIDと等価かを判定する
func (u UserID) Equals(other UserID) bool {
	return u.value == other.value
}

// IsEmpty は空のUserIDかを判定する
func (u UserID) IsEmpty() bool {
	return u.value == uuid.Nil
}