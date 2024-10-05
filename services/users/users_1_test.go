package services_users

import (
	"backend/models"
	repositories_users "backend/repositories/users"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_FetchUsers(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// モックの挙動を設定
	mockUsers := []models.UserData{
		{ID: "1", Name: "John Doe", Email: "john@example.com"},
		{ID: "2", Name: "Jane Doe", Email: "jane@example.com"},
	}
	mockUserRepository.On("FetchUsers").Return(mockUsers, nil)

	// サービス層メソッドの実行
	users, err := userService.FetchUsers()

	// エラーチェック
	assert.NoError(t, err)

	// データが期待通りか確認
	assert.Len(t, users, 2)
	assert.Equal(t, "John Doe", users[0].Name)
	assert.Equal(t, "Jane Doe", users[1].Name)

	// モックが期待通りに呼び出されたかを確認
	mockUserRepository.AssertExpectations(t)
}

func TestService_FetchUsers_EmptyList(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// 2. ユーザーが存在しない場合
	mockUserRepository.On("FetchUsers").Return([]models.UserData{}, nil)

	users, err := userService.FetchUsers()

	// エラーチェック
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 0)

	// モックが期待通りに呼び出されたかを確認
	mockUserRepository.AssertExpectations(t)
}

func TestService_FetchUserByEmailAndPassword(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// モックの挙動を設定
	mockUser := &models.UserData{
		ID:    "1",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockUserRepository.On("FetchUserByEmailAndPassword", "john@example.com", "password123").Return(mockUser, nil)

	// サービス層メソッドの実行
	user, err := userService.FetchUserByEmailAndPassword("john@example.com", "password123")

	// エラーチェック
	assert.NoError(t, err)

	// データが期待通りか確認
	assert.NotNil(t, user)                 // userがnilでないことを確認
	assert.Equal(t, "John Doe", user.Name) // ユーザーの名前が期待通りかを確認

	// モックが期待通りに呼び出されたかを確認
	mockUserRepository.AssertExpectations(t)
}

func TestService_FetchUserByEmailAndPassword_InvalidCases(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// 1. メールアドレスとパスワードが空の場合
	_, err := userService.FetchUserByEmailAndPassword("", "")
	assert.Error(t, err)
	assert.Equal(t, "email and password are required", err.Error())

	// 2. メールアドレスの形式が無効な場合
	_, err = userService.FetchUserByEmailAndPassword("invalid-email", "password123")
	assert.Error(t, err)
	assert.Equal(t, "invalid email format", err.Error())

	// 3. ユーザーが見つからない場合
	mockUserRepository.On("FetchUserByEmailAndPassword", "john@example.com", "password123").Return(nil, sql.ErrNoRows)

	_, err = userService.FetchUserByEmailAndPassword("john@example.com", "password123")
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())

	// モックが期待通りに呼び出されたかを確認
	mockUserRepository.AssertExpectations(t)
}

func TestService_FetchUserById(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// 正常系: ユーザーが見つかる場合
	expectedUser := &models.UserData{ID: "1", Name: "John Doe", Email: "john@example.com"}
	mockUserRepository.On("FetchUserById", "1").Return(expectedUser, nil)

	user, err := userService.FetchUserById("1")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)

	// モックが期待通りに呼び出されたかを確認
	mockUserRepository.AssertExpectations(t)
}

func TestService_FetchUserById_InvalidCases(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// 異常系: ユーザーが見つからない場合
	mockUserRepository.On("FetchUserById", "2").Return(nil, sql.ErrNoRows)

	user, err := userService.FetchUserById("2")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, sql.ErrNoRows, err)

	// モックが期待通りに呼び出されたかを確認
	mockUserRepository.AssertExpectations(t)
}

func TestService_FetchUserByEmail(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// 正常系: ユーザーが見つかる場合
	expectedUser := &models.UserData{ID: "1", Name: "John Doe", Email: "john@example.com"}
	mockUserRepository.On("FetchUserByEmail", "john@example.com").Return(expectedUser, nil)

	user, err := userService.FetchUserByEmail("john@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)

	// モックが期待通りに呼び出されたかを確認
	mockUserRepository.AssertExpectations(t)
}

func TestService_FetchUserByEmail_InvalidCases(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// 異常系: ユーザーが見つからない場合
	mockUserRepository.On("FetchUserByEmail", "unknown@example.com").Return(nil, sql.ErrNoRows)

	user, err := userService.FetchUserByEmail("unknown@example.com")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, sql.ErrNoRows, err)

	// モックが期待通りに呼び出されたかを確認
	mockUserRepository.AssertExpectations(t)
}
