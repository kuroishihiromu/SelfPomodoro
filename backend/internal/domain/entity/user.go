package entity

import (
	"time"

	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
)

// User はユーザーを表すドメインエンティティ
type User struct {
	ID         userVO.UserID       // Cognito sub
	Name       userVO.UserName     // 表示名
	Email      userVO.EmailAddress // メールアドレス
	Provider   userVO.Provider     // "Cognito_UserPool", "Google"
	ProviderID *string             // Google sub等
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// UserCreationParams はユーザー作成時のパラメータ
type UserCreationParams struct {
	UserID     userVO.UserID
	Name       userVO.UserName
	Email      userVO.EmailAddress
	Provider   userVO.Provider
	ProviderID *string
}

// ドメインルール：ユーザー作成ファクトリーメソッド

// NewUser は新しいユーザーを作成する
func NewUser(params UserCreationParams) *User {
	now := time.Now()
	return &User{
		ID:         params.UserID,
		Name:       params.Name,
		Email:      params.Email,
		Provider:   params.Provider,
		ProviderID: params.ProviderID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// NewUserFromStrings は文字列パラメータから新しいユーザーを作成する
func NewUserFromStrings(userID userVO.UserID, name, email, provider string, providerID *string) (*User, error) {
	userName, err := userVO.NewUserName(name)
	if err != nil {
		return nil, err
	}
	
	emailAddress, err := userVO.NewEmailAddress(email)
	if err != nil {
		return nil, err
	}
	
	userProvider, err := userVO.NewProvider(provider)
	if err != nil {
		return nil, err
	}
	
	params := UserCreationParams{
		UserID:     userID,
		Name:       userName,
		Email:      emailAddress,
		Provider:   userProvider,
		ProviderID: providerID,
	}
	
	return NewUser(params), nil
}

// ドメインルール：ユーザー検証・状態管理

// IsGoogleUser はGoogleユーザーかどうかを返す
func (u *User) IsGoogleUser() bool {
	return u.Provider.IsGoogle()
}

// IsCognitoUser はCognitoユーザーかどうかを返す
func (u *User) IsCognitoUser() bool {
	return u.Provider.IsCognito()
}

// ValidateEmail はメールアドレスの基本的な検証を行う
func (u *User) ValidateEmail() bool {
	return !u.Email.IsEmpty()
}

// ValidateName は名前の基本的な検証を行う
func (u *User) ValidateName() bool {
	return u.Name.IsValid()
}

// IsValidForCreation はユーザー作成時の必須項目を検証する
func (u *User) IsValidForCreation() bool {
	return !u.ID.IsEmpty() &&
		u.ValidateEmail() &&
		u.ValidateName() &&
		!u.Provider.IsEmpty()
}

// UpdateProfile はユーザープロフィールを更新する
func (u *User) UpdateProfile(name, email string) error {
	if name != "" {
		newName, err := userVO.NewUserName(name)
		if err != nil {
			return err
		}
		u.Name = newName
	}
	if email != "" {
		newEmail, err := userVO.NewEmailAddress(email)
		if err != nil {
			return err
		}
		u.Email = newEmail
	}
	u.UpdatedAt = time.Now()
	return nil
}

// ドメインルール：プロバイダー別ビジネスロジック

// GetProviderDisplayName はプロバイダーの表示名を返す
func (u *User) GetProviderDisplayName() string {
	return u.Provider.GetDisplayName()
}

// CanChangeEmail はメールアドレス変更可能かを判定する
func (u *User) CanChangeEmail() bool {
	return u.Provider.CanChangeEmail()
}

// GetAvatarURL はアバターURLを生成する（将来の拡張用）
func (u *User) GetAvatarURL() string {
	if u.IsGoogleUser() {
		// Google ユーザーの場合はGravatarやGoogle Photos APIを使用
		// 現在は簡易実装
		return ""
	}
	// Cognito ユーザーの場合はGravatarを使用
	return ""
}


