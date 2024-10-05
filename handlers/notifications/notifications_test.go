package handlers_notifications

import (
	"backend/models"
	services_notifications "backend/services/notifications"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetNotifications(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/notifications", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスのインスタンス化
	mockNotificationService := new(services_notifications.MockNotificationService)
	// NotificationHandlerのインスタンス化
	handler := NewNotificationHandler(mockNotificationService)

	// モックデータの設定
	mockNotifications := []models.NotificationData{
		{ID: "1", UserId: "user1", Message: "New reservation confirmed"},
		{ID: "2", UserId: "user2", Message: "Reservation canceled"},
	}
	mockNotificationService.On("FetchNotifications").Return(mockNotifications, nil)

	// ハンドラーを実行
	handler.GetNotifications(c)

	// ステータスコードとデータの確認
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "New reservation confirmed")
	assert.Contains(t, rec.Body.String(), "Reservation canceled")

	// モックが期待通りに呼び出されたかを確認
	mockNotificationService.AssertExpectations(t)
}

func TestHandler_GetNotifications_DatabaseError(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/notifications", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスのインスタンス化
	mockNotificationService := new(services_notifications.MockNotificationService)
	// NotificationHandlerのインスタンス化
	handler := NewNotificationHandler(mockNotificationService)

	// モックデータの設定
	mockNotificationService.On("FetchNotifications").Return(nil, errors.New("database error"))

	// ハンドラーを実行
	handler.GetNotifications(c)

	// ステータスコードとデータの確認
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to fetch notifications")

	// モックが期待通りに呼び出されたかを確認
	mockNotificationService.AssertExpectations(t)
}

func TestHandler_GetNotifications_NoDatas(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/notifications", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスのインスタンス化
	mockNotificationService := new(services_notifications.MockNotificationService)
	// NotificationHandlerのインスタンス化
	handler := NewNotificationHandler(mockNotificationService)

	// モックデータの設定
	mockNotificationService.On("FetchNotifications").Return([]models.NotificationData{}, nil)

	// ハンドラーを実行
	handler.GetNotifications(c)

	// ステータスコードとデータの確認
	assert.Equal(t, http.StatusOK, rec.Code)
	// 空の配列が返されることを確認
	assert.JSONEq(t, "[]", rec.Body.String())

	// モックが期待通りに呼び出されたかを確認
	mockNotificationService.AssertExpectations(t)
}

func TestHandler_AddNotification(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"user_id":"user1", "reservation_id":"reservation1", "message":"Reservation confirmed"}`
	req := httptest.NewRequest(http.MethodPost, "/api/notification", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをモック
	mockNotificationService := new(services_notifications.MockNotificationService)

	// ハンドラのインスタンス化
	handler := NewNotificationHandler(mockNotificationService)

	// モックデータの設定
	mockNotificationService.On("CreateNotification", "user1", "reservation1", "Reservation confirmed").Return(nil)

	// ハンドラーを実行
	handler.AddNotification(c)

	// ステータスコードとデータの確認
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "Notification created successfully")

	// モックが期待通りに呼び出されたかを確認
	mockNotificationService.AssertExpectations(t)
}

func TestHandler_AddNotification_ValidateError(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"user_id":"", "reservation_id":"", "message":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/notification", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをモック
	mockNotificationService := new(services_notifications.MockNotificationService)
	// ハンドラのインスタンス化
	handler := NewNotificationHandler(mockNotificationService)

	// モックデータの設定
	mockNotificationService.On("CreateNotification", "", "", "").Return(errors.New("userID, ReservationID, and message are required"))

	// ハンドラーを実行
	handler.AddNotification(c)

	// ステータスコードとデータの確認
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "UserID, ReservationID, and message are required")

	// モックが期待通りに呼び出されたかを確認
	mockNotificationService.AssertExpectations(t)
}

func TestHandler_AddNotification_UserNotFound(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"user_id":"user1", "reservation_id":"reservation1", "message":"Reservation confirmed"}`
	req := httptest.NewRequest(http.MethodPost, "/api/notification", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをモック
	mockNotificationService := new(services_notifications.MockNotificationService)

	// ハンドラのインスタンス化
	handler := NewNotificationHandler(mockNotificationService)

	// モックデータの設定
	mockNotificationService.On("CreateNotification", "user1", "reservation1", "Reservation confirmed").Return(errors.New("user not found"))

	// ハンドラーを実行
	handler.AddNotification(c)

	// ステータスコードとデータの確認
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "User not found")

	// モックが期待通りに呼び出されたかを確認
	mockNotificationService.AssertExpectations(t)
}

func TestHandler_AddNotification_ReservationNotFound(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"user_id":"user1", "reservation_id":"reservation1", "message":"Reservation confirmed"}`
	req := httptest.NewRequest(http.MethodPost, "/api/notification", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをモック
	mockNotificationService := new(services_notifications.MockNotificationService)

	// ハンドラのインスタンス化
	handler := NewNotificationHandler(mockNotificationService)

	// モックデータの設定
	mockNotificationService.On("CreateNotification", "user1", "reservation1", "Reservation confirmed").Return(errors.New("reservation not found"))

	// ハンドラーを実行
	handler.AddNotification(c)

	// ステータスコードとデータの確認
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Reservation not found")

	// モックが期待通りに呼び出されたかを確認
	mockNotificationService.AssertExpectations(t)
}

func TestHandler_AddNotification_CreateError(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"user_id":"user1", "reservation_id":"reservation1", "message":"Reservation confirmed"}`
	req := httptest.NewRequest(http.MethodPost, "/api/notification", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをモック
	mockNotificationService := new(services_notifications.MockNotificationService)

	// ハンドラのインスタンス化
	handler := NewNotificationHandler(mockNotificationService)

	// モックデータの設定
	mockNotificationService.On("CreateNotification", "user1", "reservation1", "Reservation confirmed").Return(errors.New("failed to create notification"))

	// ハンドラーを実行
	handler.AddNotification(c)

	// ステータスコードとデータの確認
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to create notification")

	// モックが期待通りに呼び出されたかを確認
	mockNotificationService.AssertExpectations(t)
}

func TestHandler_AddNotification_ServerError(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"user_id":"user1", "reservation_id":"reservation1", "message":"Reservation confirmed"}`
	req := httptest.NewRequest(http.MethodPost, "/api/notification", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをモック
	mockNotificationService := new(services_notifications.MockNotificationService)

	// ハンドラのインスタンス化
	handler := NewNotificationHandler(mockNotificationService)

	// モックデータの設定
	mockNotificationService.On("CreateNotification", "user1", "reservation1", "Reservation confirmed").Return(errors.New("server error"))

	// ハンドラーを実行
	handler.AddNotification(c)

	// ステータスコードとデータの確認
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to create reservation")

	// モックが期待通りに呼び出されたかを確認
	mockNotificationService.AssertExpectations(t)
}
