package handlers_reservations

import (
	"backend/auth"
	"backend/models"
	services_notifications "backend/services/notifications"
	services_reservations "backend/services/reservations"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// テストコード
func TestGetReservations(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/reservations", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := new(services_reservations.MockReservationService)
	handler := NewReservationHandler(nil, mockService, nil)

	// モックデータの設定
	reservationDate1, _ := time.Parse(time.RFC3339, "2024-10-01T18:00:00Z")
	reservationDate2, _ := time.Parse(time.RFC3339, "2024-10-02T19:00:00Z")

	mockReservations := []models.ReservationData{
		{ID: "1", UserId: "user1", ReservationDate: reservationDate1, NumPeople: 2, Status: "confirmed"},
		{ID: "2", UserId: "user2", ReservationDate: reservationDate2, NumPeople: 4, Status: "pending"},
	}
	mockService.On("FetchReservations").Return(mockReservations, nil)

	// ハンドラーを実行
	if assert.NoError(t, handler.GetReservations(c)) {
		// ステータスコードの確認
		assert.Equal(t, http.StatusOK, rec.Code)

		// レスポンス内容の確認
		assert.Contains(t, rec.Body.String(), "2024-10-01T18:00:00Z")
		assert.Contains(t, rec.Body.String(), "2024-10-02T19:00:00Z")
	}

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}

func TestAddReservation(t *testing.T) {
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
	if assert.NoError(t, handler.AddReservation(c)) {
		// ステータスコードの確認
		assert.Equal(t, http.StatusCreated, rec.Code)

		// レスポンス内容の確認
		assert.Contains(t, rec.Body.String(), "Reservation created successfully")
	}

	// モックが期待通りに呼び出されたかを確認
	mockReservationService.AssertExpectations(t)
	mockNotificationService.AssertExpectations(t)
}

func TestAddReservation_UserNotFound(t *testing.T) {
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
	if assert.NoError(t, handler.AddReservation(c)) {
		// ステータスコードの確認
		assert.Equal(t, http.StatusNotFound, rec.Code)

		// レスポンス内容の確認
		assert.Contains(t, rec.Body.String(), "User not found")
	}

	// モックが期待通りに呼び出されたかを確認
	mockReservationService.AssertExpectations(t)
}
