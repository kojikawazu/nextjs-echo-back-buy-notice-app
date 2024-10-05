package services_reservations

import (
	"backend/models"
	repositories_reservations "backend/repositories/reservations"
	repositories_users "backend/repositories/users"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_FetchReservations(t *testing.T) {
	// モックリポジトリをインスタンス化
	userRepository := new(repositories_users.MockUserRepository)
	reservationRepository := new(repositories_reservations.MockReservationRepository)
	reserationService := NewReservationService(userRepository, reservationRepository)

	// モックの挙動を設定
	mockReservations := []models.ReservationData{
		{ID: "1", NumPeople: 2, SpecialRequest: "No smoking", Status: "reserved"},
		{ID: "2", NumPeople: 4, SpecialRequest: "Window seat", Status: "reserved"},
	}
	reservationRepository.On("FetchReservations").Return(mockReservations, nil)

	// サービス層メソッドの実行
	reservations, err := reserationService.FetchReservations()

	// エラーチェック
	assert.NoError(t, err)

	// データが期待通りか確認
	assert.Len(t, reservations, 2)
	assert.Equal(t, "No smoking", reservations[0].SpecialRequest)
	assert.Equal(t, "Window seat", reservations[1].SpecialRequest)

	// モックが期待通りに呼び出されたかを確認
	reservationRepository.AssertExpectations(t)
}

func TestService_FetchReservations_EmptyList(t *testing.T) {
	// モックリポジトリをインスタンス化
	userRepository := new(repositories_users.MockUserRepository)
	reservationRepository := new(repositories_reservations.MockReservationRepository)
	reserationService := NewReservationService(userRepository, reservationRepository)

	// モックの挙動を設定
	reservationRepository.On("FetchReservations").Return([]models.ReservationData{}, nil)

	// サービス層メソッドの実行
	reservations, err := reserationService.FetchReservations()

	// エラーチェック
	assert.NoError(t, err)

	// データが期待通りか確認
	assert.NotNil(t, reservations)
	assert.Len(t, reservations, 0)

	// モックが期待通りに呼び出されたかを確認
	reservationRepository.AssertExpectations(t)
}

func TestService_FetchReservationById(t *testing.T) {
	// モックリポジトリをインスタンス化
	userRepository := new(repositories_users.MockUserRepository)
	reservationRepository := new(repositories_reservations.MockReservationRepository)
	reserationService := NewReservationService(userRepository, reservationRepository)

	// モックの挙動を設定
	mockReservation := &models.ReservationData{
		ID:             "1",
		NumPeople:      2,
		SpecialRequest: "No smoking",
		Status:         "reserved",
	}
	reservationRepository.On("FetchReservationById", "1").Return(mockReservation, nil)

	// サービス層メソッドの実行
	reservation, err := reserationService.FetchReservationById("1")

	// エラーチェック
	assert.NoError(t, err)

	// データが期待通りか確認
	assert.NotNil(t, reservation)
	assert.Equal(t, "No smoking", reservation.SpecialRequest)

	// モックが期待通りに呼び出されたかを確認
	reservationRepository.AssertExpectations(t)
}

func TestService_FetchReservationById_EmptyData(t *testing.T) {
	// モックリポジトリをインスタンス化
	userRepository := new(repositories_users.MockUserRepository)
	reservationRepository := new(repositories_reservations.MockReservationRepository)
	reserationService := NewReservationService(userRepository, reservationRepository)

	// モックの挙動を設定
	reservationRepository.On("FetchReservationById", "1").Return(nil, errors.New("record not found"))

	// サービス層メソッドの実行
	reservation, err := reserationService.FetchReservationById("1")

	// エラーチェック
	assert.Error(t, err)
	// データが期待通りか確認
	assert.Nil(t, reservation)

	// モックが期待通りに呼び出されたかを確認
	reservationRepository.AssertExpectations(t)
}

func TestService_FetchReservationByUserId(t *testing.T) {
	// モックリポジトリをインスタンス化
	userRepository := new(repositories_users.MockUserRepository)
	reservationRepository := new(repositories_reservations.MockReservationRepository)
	reserationService := NewReservationService(userRepository, reservationRepository)

	// モックの挙動を設定
	mockReservation := &models.ReservationData{
		ID:             "1",
		UserId:         "1",
		NumPeople:      2,
		SpecialRequest: "No smoking",
		Status:         "reserved",
	}
	reservationRepository.On("FetchReservationByUserId", "1").Return(mockReservation, nil)

	// サービス層メソッドの実行
	reservation, err := reserationService.FetchReservationByUserId("1")

	// エラーチェック
	assert.NoError(t, err)

	// データが期待通りか確認
	assert.NotNil(t, reservation)
	assert.Equal(t, "No smoking", reservation.SpecialRequest)

	// モックが期待通りに呼び出されたかを確認
	reservationRepository.AssertExpectations(t)
}

func TestService_FetchReservationByUserId_ValidationUserI(t *testing.T) {
	// モックリポジトリをインスタンス化
	userRepository := new(repositories_users.MockUserRepository)
	reservationRepository := new(repositories_reservations.MockReservationRepository)
	reserationService := NewReservationService(userRepository, reservationRepository)

	// サービス層メソッドの実行
	reservation, err := reserationService.FetchReservationByUserId("")

	// エラーチェック
	assert.Error(t, err)
	// データが期待通りか確認
	assert.Nil(t, reservation)
	assert.Equal(t, "userId is required", err.Error())

	// モックが期待通りに呼び出されたかを確認
	reservationRepository.AssertExpectations(t)
}

func TestService_FetchReservationByUserId_EmptyData(t *testing.T) {
	// モックリポジトリをインスタンス化
	userRepository := new(repositories_users.MockUserRepository)
	reservationRepository := new(repositories_reservations.MockReservationRepository)
	reserationService := NewReservationService(userRepository, reservationRepository)

	// モックの挙動を設定
	reservationRepository.On("FetchReservationByUserId", "1").Return(nil, errors.New("reservation not found"))

	// サービス層メソッドの実行
	reservation, err := reserationService.FetchReservationByUserId("1")

	// エラーチェック
	assert.Error(t, err)
	// データが期待通りか確認
	assert.Nil(t, reservation)
	assert.Equal(t, "reservation not found", err.Error())

	// モックが期待通りに呼び出されたかを確認
	reservationRepository.AssertExpectations(t)
}
