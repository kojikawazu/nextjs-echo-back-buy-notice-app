package services_users

import (
	"backend/models"
	repositories_users "backend/repositories/users"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_CreateUser(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// モックの挙動を設定
	mockUserRepository.On("FetchUserByEmail", "john@example.com").Return(nil, nil) // ユーザーが存在しない場合
	mockUserRepository.On("CreateUser", "John Doe", "john@example.com", "password123").Return(nil)

	// サービス層メソッドの実行
	err := userService.CreateUser("John Doe", "john@example.com", "password123")

	// エラーチェック
	assert.NoError(t, err)

	// モックが期待通りに呼び出されたかを確認
	mockUserRepository.AssertExpectations(t)
}

func TestService_CreateUser_InvalidCases(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// 1. 既にユーザーが存在する場合
	existingUser := &models.UserData{ID: "1", Name: "John Doe", Email: "john@example.com"}
	mockUserRepository.On("FetchUserByEmail", "john@example.com").Return(existingUser, nil)

	err := userService.CreateUser("John Doe", "john@example.com", "password123")
	assert.Error(t, err)
	assert.Equal(t, "user already exists", err.Error())

	// 2. 名前、メール、パスワードが空の場合
	err = userService.CreateUser("", "john@example.com", "password123")
	assert.Error(t, err)
	assert.Equal(t, "name, email and password are required", err.Error())

	err = userService.CreateUser("John Doe", "", "password123")
	assert.Error(t, err)
	assert.Equal(t, "name, email and password are required", err.Error())

	err = userService.CreateUser("John Doe", "john@example.com", "")
	assert.Error(t, err)
	assert.Equal(t, "name, email and password are required", err.Error())

	// 3. ユーザー追加に失敗する場合
	mockUserRepository.On("FetchUserByEmail", "newuser@example.com").Return(nil, nil)
	mockUserRepository.On("CreateUser", "New User", "newuser@example.com", "password123").Return(errors.New("insert failed"))

	err = userService.CreateUser("New User", "newuser@example.com", "password123")
	assert.Error(t, err)
	assert.Equal(t, "failed to create user", err.Error())

	// モックが期待通りに呼び出されたかを確認
	mockUserRepository.AssertExpectations(t)
}
