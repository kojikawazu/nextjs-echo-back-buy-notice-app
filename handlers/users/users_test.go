package handlers_users

import (
	"backend/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := &MockUserService{}
	handler := NewUserHandler(mockService)

	// モックの挙動を設定
	mockUsers := []models.UserData{
		{ID: "1", Name: "John Doe", Email: "john@example.com"},
		{ID: "2", Name: "Jane Doe", Email: "jane@example.com"},
	}
	mockService.On("FetchUsers").Return(mockUsers, nil)

	// ハンドラーを実行
	if assert.NoError(t, handler.GetUsers(c)) {
		// ステータスコードの確認
		assert.Equal(t, http.StatusOK, rec.Code)

		// レスポンス内容の確認
		assert.Contains(t, rec.Body.String(), "John Doe")
		assert.Contains(t, rec.Body.String(), "Jane Doe")
	}
}

func TestGetUserByEmailAndPassword(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"email":"john@example.com", "password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := &MockUserService{}
	handler := NewUserHandler(mockService)

	// モックデータの設定
	mockUser := &models.UserData{
		ID:    "1",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockService.On("FetchUserByEmailAndPassword", "john@example.com", "password123").Return(mockUser, nil)

	// ハンドラーを実行
	if assert.NoError(t, handler.GetUserByEmailAndPassword(c)) {
		// ステータスコードの確認
		assert.Equal(t, http.StatusOK, rec.Code)

		// レスポンス内容の確認
		assert.Contains(t, rec.Body.String(), "John Doe")
	}

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}

func TestAddUser(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"name":"John Doe", "email":"john@example.com", "password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/add", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := &MockUserService{}
	handler := NewUserHandler(mockService)

	// サービス側でのモックの挙動を設定
	mockService.On("FetchUserByEmail", "john@example.com").Return(nil, nil) // ユーザーが存在しない場合
	mockService.On("CreateUser", "John Doe", "john@example.com", "password123").Return(nil)

	// ハンドラーを実行
	if assert.NoError(t, handler.AddUser(c)) {
		// ステータスコードの確認
		assert.Equal(t, http.StatusCreated, rec.Code)

		// レスポンス内容の確認
		assert.Contains(t, rec.Body.String(), "User created successfully")
	}
}

func TestAddUserAlreadyExists(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"name":"John Doe", "email":"john@example.com", "password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/add", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := &MockUserService{}
	handler := NewUserHandler(mockService)

	// サービス側でのモックの挙動を設定
	existingUser := &models.UserData{Name: "John Doe", Email: "john@example.com"}
	mockService.On("FetchUserByEmail", "john@example.com").Return(existingUser, nil) // ユーザーが既に存在する

	// ハンドラーを実行
	if assert.NoError(t, handler.AddUser(c)) {
		// ステータスコードの確認
		assert.Equal(t, http.StatusConflict, rec.Code)

		// レスポンス内容の確認
		assert.Contains(t, rec.Body.String(), "User already exists")
	}
}
