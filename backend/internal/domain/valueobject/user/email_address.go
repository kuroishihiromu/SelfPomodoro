package user

import (
	"errors"
	"regexp"
	"strings"
)

// EmailAddress はメールアドレスを表すValue Object
type EmailAddress struct {
	value string
}

// メールアドレスの基本的な正規表現パターン
var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmailAddress は新しいEmailAddressを作成する
func NewEmailAddress(email string) (EmailAddress, error) {
	if email == "" {
		return EmailAddress{}, errors.New("メールアドレスは必須です")
	}
	
	// 前後の空白を除去
	email = strings.TrimSpace(email)
	
	// 基本的なフォーマットチェック
	if !emailPattern.MatchString(email) {
		return EmailAddress{}, errors.New("無効なメールアドレス形式です")
	}
	
	// 長さチェック
	if len(email) > 254 {
		return EmailAddress{}, errors.New("メールアドレスが長すぎます")
	}
	
	return EmailAddress{value: strings.ToLower(email)}, nil
}

// MustNewEmailAddress はパニックを起こす可能性があるEmailAddress作成（テスト用）
func MustNewEmailAddress(email string) EmailAddress {
	ea, err := NewEmailAddress(email)
	if err != nil {
		panic(err)
	}
	return ea
}

// Value は文字列の値を返す
func (e EmailAddress) Value() string {
	return e.value
}

// Equals は別のEmailAddressと等価かを判定する
func (e EmailAddress) Equals(other EmailAddress) bool {
	return e.value == other.value
}

// Domain はドメイン部分を返す
func (e EmailAddress) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// LocalPart はローカル部分（@より前）を返す
func (e EmailAddress) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

// IsGmailAddress はGmailアドレスかを判定する
func (e EmailAddress) IsGmailAddress() bool {
	return e.Domain() == "gmail.com"
}

// String は文字列表現を返す
func (e EmailAddress) String() string {
	return e.value
}

// IsEmpty はメールアドレスが空かどうかを判定する
func (e EmailAddress) IsEmpty() bool {
	return e.value == ""
}