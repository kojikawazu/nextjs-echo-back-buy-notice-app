package services_reservations

import (
	"errors"
	"testing"

	"backend/models"
	repositories_reservations "backend/repositories/reservations"
	repositories_users "backend/repositories/users"

	"github.com/stretchr/testify/assert"
)

func TestService_CreateReservation_Success(t *testing.T) {
	// モックリポジトリをインスタンス化
	userRepository := new(repositories_users.MockUserRepository)
	reservationRepository := new(repositories_reservations.MockReservationRepository)
	reserationService := NewReservationService(userRepository, reservationRepository)

	// ユーザーが存在する場合のモックの挙動を設定
	userRepository.On("FetchUserById", "user1").Return(&models.UserData{ID: "user1", Name: "John Doe", Email: "john@example.com"}, nil)

	// 予約作成のモックの挙動を設定
	reservationRepository.On("CreateReservation", "user1", "2024-10-10 12:00:00", 4, "Special request", "pending").Return("reservation1", nil)

	// サービス層メソッドの実行
	reservationId, err := reserationService.CreateReservation("user1", "2024-10-10 12:00:00", 4, "Special request", "")

	// エラーチェックと結果の確認
	assert.NoError(t, err)
	assert.Equal(t, "reservation1", reservationId)

	// モックが期待通りに呼び出されたかを確認
	userRepository.AssertExpectations(t)
	reservationRepository.AssertExpectations(t)
}

func TestService_CreateReservation_ValidationError(t *testing.T) {
	// モックリポジトリをインスタンス化
	userRepository := new(repositories_users.MockUserRepository)
	reservationRepository := new(repositories_reservations.MockReservationRepository)
	reserationService := NewReservationService(userRepository, reservationRepository)

	// バリデーションエラーを確認するため、ユーザー取得などは不要
	_, err := reserationService.CreateReservation("user1", "", 0, "Special request", "")

	// エラーチェック
	assert.Error(t, err)
	assert.Equal(t, "userID, reservation date, and num_people are required", err.Error())

	// モックが呼び出されていないか確認
	userRepository.AssertNotCalled(t, "FetchUserById")
	reservationRepository.AssertNotCalled(t, "CreateReservation")
}

func TestService_CreateReservation_InvalidDate(t *testing.T) {
	// モックリポジトリをインスタンス化
	userRepository := new(repositories_users.MockUserRepository)
	reservationRepository := new(repositories_reservations.MockReservationRepository)
	reserationService := NewReservationService(userRepository, reservationRepository)

	// ユーザーが存在する場合のモックの挙動を設定
	userRepository.On("FetchUserById", "user1").Return(&models.UserData{ID: "user1", Name: "John Doe", Email: "john@example.com"}, nil)

	// 不正な日付フォーマットを渡す
	_, err := reserationService.CreateReservation("user1", "invalid-date", 4, "Special request", "")

	// エラーチェック
	assert.Error(t, err)
	assert.Equal(t, "invalid reservation date format. Use 'YYYY-MM-DD HH:MM:SS'", err.Error())

	// モックが呼び出されていないか確認
	reservationRepository.AssertNotCalled(t, "CreateReservation")
}

func TestService_CreateReservation_UserNotFound(t *testing.T) {
	// モックリポジトリをインスタンス化
	userRepository := new(repositories_users.MockUserRepository)
	reservationRepository := new(repositories_reservations.MockReservationRepository)
	reserationService := NewReservationService(userRepository, reservationRepository)

	// ユーザーが存在しない場合のモックの挙動を設定
	userRepository.On("FetchUserById", "user1").Return(nil, errors.New("user not found"))

	// サービス層メソッドの実行
	_, err := reserationService.CreateReservation("user1", "2024-10-10 12:00:00", 4, "Special request", "")

	// エラーチェック
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())

	// モックが期待通りに呼び出されたか確認
	userRepository.AssertCalled(t, "FetchUserById", "user1")
	reservationRepository.AssertNotCalled(t, "CreateReservation")
}
