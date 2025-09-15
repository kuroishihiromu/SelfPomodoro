package user

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// UserName はユーザー名を表すValue Object
type UserName struct {
	value string
}

// NewUserName は新しいUserNameを作成する
func NewUserName(name string) (UserName, error) {
	if name == "" {
		return UserName{}, errors.New("ユーザー名は必須です")
	}
	
	// 前後の空白を除去
	name = strings.TrimSpace(name)
	
	// 長さチェック（UTF-8文字数で）
	if utf8.RuneCountInString(name) < 1 {
		return UserName{}, errors.New("ユーザー名は1文字以上で設定してください")
	}
	
	if utf8.RuneCountInString(name) > 50 {
		return UserName{}, errors.New("ユーザー名は50文字以下で設定してください")
	}
	
	// 無効な文字のチェック（制御文字など）
	for _, r := range name {
		if r < 32 && r != 9 && r != 10 && r != 13 { // 制御文字（タブ、改行、復帰以外）
			return UserName{}, errors.New("無効な文字が含まれています")
		}
	}
	
	return UserName{value: name}, nil
}

// MustNewUserName はパニックを起こす可能性があるUserName作成（テスト用）
func MustNewUserName(name string) UserName {
	un, err := NewUserName(name)
	if err != nil {
		panic(err)
	}
	return un
}

// Value は文字列の値を返す
func (u UserName) Value() string {
	return u.value
}

// Equals は別のUserNameと等価かを判定する
func (u UserName) Equals(other UserName) bool {
	return u.value == other.value
}

// Length は文字数を返す
func (u UserName) Length() int {
	return utf8.RuneCountInString(u.value)
}

// IsShort は短い名前かを判定する（5文字以下）
func (u UserName) IsShort() bool {
	return u.Length() <= 5
}

// IsLong は長い名前かを判定する（20文字以上）
func (u UserName) IsLong() bool {
	return u.Length() >= 20
}

// ContainsEmail はメールアドレスらしき文字列を含むかを判定する
func (u UserName) ContainsEmail() bool {
	return strings.Contains(u.value, "@") && strings.Contains(u.value, ".")
}

// String は文字列表現を返す
func (u UserName) String() string {
	return u.value
}

// IsEmpty はユーザー名が空かどうかを判定する
func (u UserName) IsEmpty() bool {
	return u.value == ""
}

// IsValid はユーザー名が有効かどうかを判定する
func (u UserName) IsValid() bool {
	return !u.IsEmpty() && u.Length() >= 1 && u.Length() <= 50
}