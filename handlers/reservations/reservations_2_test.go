package handlers_reservations

import (
	"backend/auth"
	services_notifications "backend/services/notifications"
	services_reservations "backend/services/reservations"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_AddReservation(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"reservation_date":"2024-10-01 18:00:00", "num_people":2, "special_request":"Window seat", "status":"confirmed"}`
	req := httptest.NewRequest(http.MethodPost, "/api/reservation", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockReservationService := new(services_reservations.MockReservationService)
	mockNotificationService := new(services_notifications.MockNotificationService)
	handler := NewReservationHandler(nil, mockReservationService, mockNotificationService)

	// JWTトークンのモックを作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Claims{
		UserID: "user1",
	})
	tokenString, _ := token.SignedString(auth.JwtKey) // 実際のキーを使用してトークンを生成

	cookie := &http.Cookie{
		Name:  "token",
		Value: tokenString,
	}
	req.AddCookie(cookie)

	// モックデータの設定
	mockReservationService.On("CreateReservation", "user1", "2024-10-01 18:00:00", 2, "Window seat", "confirmed").Return("reservationId", nil)
	mockNotificationService.On("CreateNotification", "user1", "reservationId", "New reservation created for user user1").Return(nil)

	// ハンドラーを実行
	handler.AddReservation(c)

	// ステータスコードの確認
	assert.Equal(t, http.StatusCreated, rec.Code)

	// レスポンス内容の確認
	assert.Contains(t, rec.Body.String(), "Reservation created successfully")

	// モックが期待通りに呼び出されたかを確認
	mockReservationService.AssertExpectations(t)
	mockNotificationService.AssertExpectations(t)
}

func TestHandler_AddReservation_ValidationError(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"reservation_date":"2024-10-01 18:00:00", "num_people":2, "special_request":"Window seat", "status":"confirmed"}`
	req := httptest.NewRequest(http.MethodPost, "/api/reservation", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockReservationService := new(services_reservations.MockReservationService)
	handler := NewReservationHandler(nil, mockReservationService, nil)

	// JWTトークンのモックを作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Claims{
		UserID: "user1",
	})
	tokenString, _ := token.SignedString(auth.JwtKey)

	cookie := &http.Cookie{
		Name:  "token",
		Value: tokenString,
	}
	req.AddCookie(cookie)

	// 予約作成時にモックを設定（通常はここでエラーが返るが、ユーザーが存在しないため不要）
	mockReservationService.On("CreateReservation", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return("", errors.New("userID, reservation date, and num_people are required"))

	// ハンドラーを実行
	handler.AddReservation(c)

	// ステータスコードの確認
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// レスポンス内容の確認
	assert.Contains(t, rec.Body.String(), "UserID, reservation date, and num_people are required")

	// モックが期待通りに呼び出されたかを確認
	mockReservationService.AssertExpectations(t)
}

func TestHandler_AddReservation_ValidationDateFormat(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"reservation_date":"2024-10-01 18:00:00", "num_people":2, "special_request":"Window seat", "status":"confirmed"}`
	req := httptest.NewRequest(http.MethodPost, "/api/reservation", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockReservationService := new(services_reservations.MockReservationService)
	handler := NewReservationHandler(nil, mockReservationService, nil)

	// JWTトークンのモックを作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Claims{
		UserID: "user1",
	})
	tokenString, _ := token.SignedString(auth.JwtKey)

	cookie := &http.Cookie{
		Name:  "token",
		Value: tokenString,
	}
	req.AddCookie(cookie)

	// 予約作成時にモックを設定（通常はここでエラーが返るが、ユーザーが存在しないため不要）
	mockReservationService.On("CreateReservation", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return("", errors.New("invalid reservation date format. Use 'YYYY-MM-DD HH:MM:SS'"))

	// ハンドラーを実行
	handler.AddReservation(c)

	// ステータスコードの確認
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// レスポンス内容の確認
	assert.Contains(t, rec.Body.String(), "Invalid reservation date format. Use 'YYYY-MM-DD HH:MM:SS'")

	// モックが期待通りに呼び出されたかを確認
	mockReservationService.AssertExpectations(t)
}

func TestHandler_AddReservation_UserNotFound(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"reservation_date":"2024-10-01 18:00:00", "num_people":2, "special_request":"Window seat", "status":"confirmed"}`
	req := httptest.NewRequest(http.MethodPost, "/api/reservation", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockReservationService := new(services_reservations.MockReservationService)
	handler := NewReservationHandler(nil, mockReservationService, nil)

	// JWTトークンのモックを作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Claims{
		UserID: "user1",
	})
	tokenString, _ := token.SignedString(auth.JwtKey)

	cookie := &http.Cookie{
		Name:  "token",
		Value: tokenString,
	}
	req.AddCookie(cookie)

	// 予約作成時にモックを設定（通常はここでエラーが返るが、ユーザーが存在しないため不要）
	mockReservationService.On("CreateReservation", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return("", errors.New("user not found")) // ここではエラーが発生することはないが、あくまで安全のため

	// ハンドラーを実行
	handler.AddReservation(c)

	// ステータスコードの確認
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// レスポンス内容の確認
	assert.Contains(t, rec.Body.String(), "User not found")

	// モックが期待通りに呼び出されたかを確認
	mockReservationService.AssertExpectations(t)
}

func TestHandler_AddReservation_CreateError(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"reservation_date":"2024-10-01 18:00:00", "num_people":2, "special_request":"Window seat", "status":"confirmed"}`
	req := httptest.NewRequest(http.MethodPost, "/api/reservation", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockReservationService := new(services_reservations.MockReservationService)
	handler := NewReservationHandler(nil, mockReservationService, nil)

	// JWTトークンのモックを作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Claims{
		UserID: "user1",
	})
	tokenString, _ := token.SignedString(auth.JwtKey)

	cookie := &http.Cookie{
		Name:  "token",
		Value: tokenString,
	}
	req.AddCookie(cookie)

	// 予約作成時にモックを設定（通常はここでエラーが返るが、ユーザーが存在しないため不要）
	mockReservationService.On("CreateReservation", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return("", errors.New("failed to create reservation"))

	// ハンドラーを実行
	handler.AddReservation(c)

	// ステータスコードの確認
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// レスポンス内容の確認
	assert.Contains(t, rec.Body.String(), "Failed to create reservation")

	// モックが期待通りに呼び出されたかを確認
	mockReservationService.AssertExpectations(t)
}

func TestHandler_AddReservation_ServerError(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"reservation_date":"2024-10-01 18:00:00", "num_people":2, "special_request":"Window seat", "status":"confirmed"}`
	req := httptest.NewRequest(http.MethodPost, "/api/reservation", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockReservationService := new(services_reservations.MockReservationService)
	handler := NewReservationHandler(nil, mockReservationService, nil)

	// JWTトークンのモックを作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Claims{
		UserID: "user1",
	})
	tokenString, _ := token.SignedString(auth.JwtKey)

	cookie := &http.Cookie{
		Name:  "token",
		Value: tokenString,
	}
	req.AddCookie(cookie)

	// 予約作成時にモックを設定（通常はここでエラーが返るが、ユーザーが存在しないため不要）
	mockReservationService.On("CreateReservation", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return("", errors.New("server error"))

	// ハンドラーを実行
	handler.AddReservation(c)

	// ステータスコードの確認
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// レスポンス内容の確認
	assert.Contains(t, rec.Body.String(), "Failed to create reservation")

	// モックが期待通りに呼び出されたかを確認
	mockReservationService.AssertExpectations(t)
}
