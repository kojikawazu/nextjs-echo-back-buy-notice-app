package handlers_reservations

import (
	"backend/models"
	services_reservations "backend/services/reservations"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// テストコード
func TestHandler_GetReservations(t *testing.T) {
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
	handler.GetReservations(c)

	// ステータスコードの確認
	assert.Equal(t, http.StatusOK, rec.Code)

	// レスポンス内容の確認
	assert.Contains(t, rec.Body.String(), "2024-10-01T18:00:00Z")
	assert.Contains(t, rec.Body.String(), "2024-10-02T19:00:00Z")

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}

func TestHandler_GetReservations_ServiceError(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/reservations", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := new(services_reservations.MockReservationService)
	handler := NewReservationHandler(nil, mockService, nil)

	// モックデータの設定

	mockService.On("FetchReservations").Return(nil, errors.New("database error"))

	// ハンドラーを実行
	handler.GetReservations(c)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to fetch reservations")

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}

func TestHandler_GetReservations_NoData(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/reservations", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := new(services_reservations.MockReservationService)
	handler := NewReservationHandler(nil, mockService, nil)

	// モックデータの設定

	mockService.On("FetchReservations").Return([]models.ReservationData{}, nil)

	// ハンドラーを実行
	handler.GetReservations(c)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, "[]", rec.Body.String())

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}

func TestHandler_GetReservationByUserId(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/reservations/:user_id", nil)
	rec := httptest.NewRecorder()

	// コンテキストにユーザーIDを設定
	c := e.NewContext(req, rec)
	c.SetParamNames("user_id")
	c.SetParamValues("user1")

	// モックサービスをインスタンス化
	mockService := new(services_reservations.MockReservationService)
	handler := NewReservationHandler(nil, mockService, nil)

	// モックの挙動を設定
	mockReservation := &models.ReservationData{
		ID:             "reservation1",
		UserId:         "user1",
		SpecialRequest: "Reservation 1",
		CreatedAt:      time.Now(),
	}
	mockService.On("FetchReservationByUserId", "user1").Return(mockReservation, nil)

	// ハンドラーを実行
	handler.GetReservationByUserId(c)

	// ステータスコードとレスポンス内容を確認
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Reservation 1")

	// モックが期待通りに呼び出されたか確認
	mockService.AssertExpectations(t)
}

func TestHandler_GetReservationByUserId_ValidationError(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/reservations/:user_id", nil)
	rec := httptest.NewRecorder()

	// コンテキストにユーザーIDを設定
	c := e.NewContext(req, rec)
	c.SetParamNames("user_id")
	c.SetParamValues("user1")

	// モックサービスをインスタンス化
	mockService := new(services_reservations.MockReservationService)
	handler := NewReservationHandler(nil, mockService, nil)

	// モックの挙動を設定
	mockService.On("FetchReservationByUserId", mock.Anything).Return(nil, errors.New("userId is required"))

	// ハンドラーを実行
	handler.GetReservationByUserId(c)

	// ステータスコードとレスポンス内容を確認
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "UserId is required")

	// モックが期待通りに呼び出されたか確認
	mockService.AssertExpectations(t)
}

func TestHandler_GetReservationByUserId_NoReservationData(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/reservations/:user_id", nil)
	rec := httptest.NewRecorder()

	// コンテキストにユーザーIDを設定
	c := e.NewContext(req, rec)
	c.SetParamNames("user_id")
	c.SetParamValues("user1")

	// モックサービスをインスタンス化
	mockService := new(services_reservations.MockReservationService)
	handler := NewReservationHandler(nil, mockService, nil)

	// モックの挙動を設定
	mockService.On("FetchReservationByUserId", mock.Anything).Return(nil, errors.New("reservation not found"))

	// ハンドラーを実行
	handler.GetReservationByUserId(c)

	// ステータスコードとレスポンス内容を確認
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Reservation not found")

	// モックが期待通りに呼び出されたか確認
	mockService.AssertExpectations(t)
}

func TestHandler_GetReservationByUserId_ServerError(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/reservations/:user_id", nil)
	rec := httptest.NewRecorder()

	// コンテキストにユーザーIDを設定
	c := e.NewContext(req, rec)
	c.SetParamNames("user_id")
	c.SetParamValues("user1")

	// モックサービスをインスタンス化
	mockService := new(services_reservations.MockReservationService)
	handler := NewReservationHandler(nil, mockService, nil)

	// モックの挙動を設定
	mockService.On("FetchReservationByUserId", mock.Anything).Return(nil, errors.New("server error"))

	// ハンドラーを実行
	handler.GetReservationByUserId(c)

	// ステータスコードとレスポンス内容を確認
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to not reservation")

	// モックが期待通りに呼び出されたか確認
	mockService.AssertExpectations(t)
}
