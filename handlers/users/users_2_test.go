package handlers_users

import (
	services_users "backend/services/users"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_AddUser(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"name":"John Doe", "email":"john@example.com", "password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/add", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := new(services_users.MockUserService)
	handler := NewUserHandler(mockService)

	// サービス側でのモックの挙動を設定
	mockService.On("CreateUser", "John Doe", "john@example.com", "password123").Return(nil)

	// ハンドラーを実行
	handler.AddUser(c)

	// ステータスコードとレスポンス内容を確認
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "User created successfully")

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}

func TestHandler_AddUser_ValidationError(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"name":"", "email":"", "password":""}` // 空のフィールド
	req := httptest.NewRequest(http.MethodPost, "/api/user/add", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := new(services_users.MockUserService)
	handler := NewUserHandler(mockService)

	// サービス側でのモックの挙動を設定
	mockService.On("CreateUser", "", "", "").Return(errors.New("name, email and password are required"))

	// ハンドラーを実行
	handler.AddUser(c)

	// ステータスコードとレスポンス内容を確認
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Name, email and password are required")

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}

func TestHandler_AddUser_InvalidEmailFormat(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"name":"John Doe", "email":"invalid-email", "password":"password123"}` // 無効なメール形式
	req := httptest.NewRequest(http.MethodPost, "/api/user/add", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := new(services_users.MockUserService)
	handler := NewUserHandler(mockService)

	// サービス側でのモックの挙動を設定
	mockService.On("CreateUser", "John Doe", "invalid-email", "password123").Return(errors.New("invalid email format"))

	// ハンドラーを実行
	handler.AddUser(c)

	// ステータスコードとレスポンス内容を確認
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid email format")

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}

func TestHandler_AddUser_AlreadyExists(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"name":"John Doe", "email":"john@example.com", "password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/add", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := new(services_users.MockUserService)
	handler := NewUserHandler(mockService)

	// サービス側でのモックの挙動を設定
	mockService.On("CreateUser", "John Doe", "john@example.com", "password123").Return(errors.New("user already exists"))

	// ハンドラーを実行
	handler.AddUser(c)

	// ステータスコードとレスポンス内容を確認
	assert.Equal(t, http.StatusConflict, rec.Code)
	assert.Contains(t, rec.Body.String(), "User already exists")

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}

func TestHandler_AddUser_ServiceError(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"name":"John Doe", "email":"john@example.com", "password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/add", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := new(services_users.MockUserService)
	handler := NewUserHandler(mockService)

	// サービス側でのモックの挙動を設定
	mockService.On("CreateUser", "John Doe", "john@example.com", "password123").Return(errors.New("failed to create user"))

	// ハンドラーを実行
	handler.AddUser(c)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to create user")

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}
