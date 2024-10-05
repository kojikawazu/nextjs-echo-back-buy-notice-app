package handlers_notifications

import (
	handlers_reservations "backend/handlers/reservations"
	"backend/models"
	repositories_notifications "backend/repositories/notifications"
	services_notifications "backend/services/notifications"
	services_users "backend/services/users"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetNotifications(t *testing.T) {
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
	if assert.NoError(t, handler.GetNotifications(c)) {
		// ステータスコードの確認
		assert.Equal(t, http.StatusOK, rec.Code)

		// レスポンス内容の確認
		assert.Contains(t, rec.Body.String(), "New reservation confirmed")
		assert.Contains(t, rec.Body.String(), "Reservation canceled")
	}

	// モックが期待通りに呼び出されたかを確認
	mockNotificationService.AssertExpectations(t)
}

func TestAddNotification(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"user_id":"user1", "reservation_id":"reservation1", "message":"Reservation confirmed"}`
	req := httptest.NewRequest(http.MethodPost, "/api/notification", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスとリポジトリをインスタンス化
	mockNotificationRepository := new(repositories_notifications.MockNotificationRepository)
	mockUserService := new(services_users.MockUserService)
	mockReservationService := new(handlers_reservations.MockReservationService)

	// NotificationServiceのインスタンス化
	notificationService := services_notifications.NewNotificationService(
		mockUserService,
		mockReservationService,
		mockNotificationRepository,
	)

	// ハンドラのインスタンス化
	handler := NewNotificationHandler(notificationService)

	// モックデータの設定
	mockUserService.On("FetchUserById", "user1").Return(&models.UserData{ID: "user1", Name: "John Doe"}, nil)
	mockReservationService.On("FetchReservationById", "reservation1").Return(&models.ReservationData{ID: "reservation1"}, nil)
	mockNotificationRepository.On("CreateNotification", "user1", "reservation1", "Reservation confirmed").Return(nil)

	// ハンドラーを実行
	if assert.NoError(t, handler.AddNotification(c)) {
		// ステータスコードの確認
		assert.Equal(t, http.StatusCreated, rec.Code)

		// レスポンス内容の確認
		assert.Contains(t, rec.Body.String(), "Notification created successfully")
	}

	// モックが期待通りに呼び出されたかを確認
	mockUserService.AssertExpectations(t)
	mockReservationService.AssertExpectations(t)
	mockNotificationRepository.AssertExpectations(t)
}

func TestAddNotification_UserNotFound(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"user_id":"user1", "reservation_id":"reservation1", "message":"Reservation confirmed"}`
	req := httptest.NewRequest(http.MethodPost, "/api/notification", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスとリポジトリをインスタンス化
	mockNotificationRepository := new(repositories_notifications.MockNotificationRepository)
	mockUserService := new(services_users.MockUserService)
	mockReservationService := new(handlers_reservations.MockReservationService)

	// NotificationServiceのインスタンス化
	notificationService := services_notifications.NewNotificationService(
		mockUserService,
		mockReservationService,
		mockNotificationRepository,
	)

	// ハンドラのインスタンス化
	handler := NewNotificationHandler(notificationService)

	// モックデータの設定
	mockUserService.On("FetchUserById", "user1").Return(nil, nil) // ユーザーが存在しない

	// ハンドラーを実行
	if assert.NoError(t, handler.AddNotification(c)) {
		// ステータスコードの確認
		assert.Equal(t, http.StatusNotFound, rec.Code)

		// レスポンス内容の確認
		assert.Contains(t, rec.Body.String(), "User not found")
	}

	// モックが期待通りに呼び出されたかを確認
	mockUserService.AssertExpectations(t)
}
