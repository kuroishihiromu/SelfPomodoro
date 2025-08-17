package mapper

import (
	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	"github.com/tsunakit99/selfpomodoro/internal/usecase/dto"
)

// UserMapper はユーザー関連のマッピングを担当する
type UserMapper struct{}

// NewUserMapper は新しいUserMapperを作成する
func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

// ToUserResponse はドメインモデルからAPIレスポンス形式に変換する
func (m *UserMapper) ToUserResponse(user *entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID.Value(),
		Name:      user.Name.Value(),
		Email:     user.Email.Value(),
		Provider:  user.Provider.Value(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}