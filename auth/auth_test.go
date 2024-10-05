package auth

import (
	"backend/models"
	services_users "backend/services/users"
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// .envファイルの読み込み
	err := godotenv.Load("../.env.test")
	if err != nil {
		panic("Error loading ../.env.test file")
	}

	// テストを実行
	code := m.Run()

	// テスト終了後に終了コードで終了
	os.Exit(code)
}

func TestLogin(t *testing.T) {
	e := echo.New()
	reqBody := `{"email":"test@example.com", "password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUserService := new(services_users.MockUserService)
	handler := NewAuthHandler(mockUserService)

	// Mockの設定
	mockUserService.On("FetchUserByEmailAndPassword", "test@example.com", "password123").Return(&models.UserData{ID: "user1", Email: "test@example.com", Name: "Test User"}, nil)

	// テスト実行
	if assert.NoError(t, handler.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Login successful")
	}

	mockUserService.AssertExpectations(t)
}

func TestLogout(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := NewAuthHandler(new(services_users.MockUserService))

	// テスト実行
	if assert.NoError(t, handler.Logout(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Logout successful")
	}
}
