package services_notifications

import (
	"backend/models"
	repositories_notifications "backend/repositories/notifications"
	repositories_reservations "backend/repositories/reservations"
	repositories_users "backend/repositories/users"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_FetchNotifications(t *testing.T) {
	// モックをインスタンス化
	notificationRepository := new(repositories_notifications.MockNotificationRepository)
	notificationService := NewNotificationService(nil, nil, notificationRepository)

	// モックの挙動を設定
	mockNotifications := []models.NotificationData{
		{ID: "1", Message: "New reservation confirmed"},
		{ID: "2", Message: "Reservation canceled"},
	}
	notificationRepository.On("FetchNotifications").Return(mockNotifications, nil)

	// サービス層メソッドの実行
	notifications, err := notificationService.FetchNotifications()

	// エラーチェック
	assert.NoError(t, err)

	// データが期待通りか確認
	assert.Len(t, notifications, 2)
	assert.Equal(t, "New reservation confirmed", notifications[0].Message)
	assert.Equal(t, "Reservation canceled", notifications[1].Message)

	// モックが期待通りに呼び出されたかを確認
	notificationRepository.AssertExpectations(t)
}

func TestService_FetchNotifications_NoDatas(t *testing.T) {
	// モックをインスタンス化
	notificationRepository := new(repositories_notifications.MockNotificationRepository)
	notificationService := NewNotificationService(nil, nil, notificationRepository)

	// モックの挙動を設定
	notificationRepository.On("FetchNotifications").Return([]models.NotificationData{}, nil)

	// サービス層メソッドの実行
	notifications, err := notificationService.FetchNotifications()

	// エラーチェック
	assert.NoError(t, err)
	assert.NotNil(t, notifications)
	assert.Len(t, notifications, 0)

	// モックが期待通りに呼び出されたかを確認
	notificationRepository.AssertExpectations(t)
}

func TestService_CreateNotification(t *testing.T) {
	// モックをインスタンス化
	userRepository := new(repositories_users.MockUserRepository)
	reservationRepository := new(repositories_reservations.MockReservationRepository)
	notificationRepository := new(repositories_notifications.MockNotificationRepository)
	notificationService := NewNotificationService(userRepository, reservationRepository, notificationRepository)

	// モックの挙動を設定
	userRepository.On("FetchUserById", "user1").Return(&models.UserData{ID: "user1", Name: "John Doe", Email: "user@example.com"}, nil)
	reservationRepository.On("FetchReservationById", "reservation1").Return(&models.ReservationData{ID: "reservation1", UserId: "user1"}, nil)
	notificationRepository.On("CreateNotification", "user1", "reservation1", "New reservation confirmed").Return(nil)

	// サービス層メソッドの実行
	err := notificationService.CreateNotification("user1", "reservation1", "New reservation confirmed")

	// エラーチェック
	assert.NoError(t, err)

	// モックが期待通りに呼び出されたかを確認
	userRepository.AssertExpectations(t)
	reservationRepository.AssertExpectations(t)
	notificationRepository.AssertExpectations(t)
}

func TestService_CreateNotification_ValidationError(t *testing.T) {
	// モックをインスタンス化
	userRepository := new(repositories_users.MockUserRepository)
	reservationRepository := new(repositories_reservations.MockReservationRepository)
	notificationRepository := new(repositories_notifications.MockNotificationRepository)
	notificationService := NewNotificationService(userRepository, reservationRepository, notificationRepository)

	// サービス層メソッドの実行
	err := notificationService.CreateNotification("", "", "")

	// エラーチェック
	assert.Error(t, err)
	assert.Equal(t, "userID, ReservationID, and message are required", err.Error())

	// モックが期待通りに呼び出されてないか確認
	userRepository.AssertNotCalled(t, "FetchUserById")
	reservationRepository.AssertNotCalled(t, "FetchReservationById")
	notificationRepository.AssertNotCalled(t, "CreateNotification")
}

func TestService_CreateNotification_NoUserData(t *testing.T) {
	// モックをインスタンス化
	userRepository := new(repositories_users.MockUserRepository)
	reservationRepository := new(repositories_reservations.MockReservationRepository)
	notificationRepository := new(repositories_notifications.MockNotificationRepository)
	notificationService := NewNotificationService(userRepository, reservationRepository, notificationRepository)

	// モックの挙動を設定
	userRepository.On("FetchUserById", "user1").Return(nil, errors.New("failed to create user"))

	// サービス層メソッドの実行
	err := notificationService.CreateNotification("user1", "reservation1", "New reservation confirmed")

	// エラーチェック
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())

	// モックが期待通りに呼び出されてないか確認
	userRepository.AssertCalled(t, "FetchUserById", "user1")
	reservationRepository.AssertNotCalled(t, "FetchReservationById")
	notificationRepository.AssertNotCalled(t, "CreateNotification")
}

func TestService_CreateNotification_NoReservationData(t *testing.T) {
	// モックをインスタンス化
	userRepository := new(repositories_users.MockUserRepository)
	reservationRepository := new(repositories_reservations.MockReservationRepository)
	notificationRepository := new(repositories_notifications.MockNotificationRepository)
	notificationService := NewNotificationService(userRepository, reservationRepository, notificationRepository)

	// モックの挙動を設定
	userRepository.On("FetchUserById", "user1").Return(&models.UserData{ID: "user1", Name: "John Doe", Email: "user@example.com"}, nil)
	reservationRepository.On("FetchReservationById", "reservation1").Return(nil, errors.New("failed to fetch reservation"))

	// サービス層メソッドの実行
	err := notificationService.CreateNotification("user1", "reservation1", "New reservation confirmed")

	// エラーチェック
	assert.Error(t, err)
	assert.Equal(t, "reservation not found", err.Error())

	// モックが期待通りに呼び出されてないか確認
	userRepository.AssertCalled(t, "FetchUserById", "user1")
	reservationRepository.AssertCalled(t, "FetchReservationById", "reservation1")
	notificationRepository.AssertNotCalled(t, "CreateNotification")
}

func TestService_CreateNotification_CreateError(t *testing.T) {
	// モックをインスタンス化
	userRepository := new(repositories_users.MockUserRepository)
	reservationRepository := new(repositories_reservations.MockReservationRepository)
	notificationRepository := new(repositories_notifications.MockNotificationRepository)
	notificationService := NewNotificationService(userRepository, reservationRepository, notificationRepository)

	// モックの挙動を設定
	userRepository.On("FetchUserById", "user1").Return(&models.UserData{ID: "user1", Name: "John Doe", Email: "user@example.com"}, nil)
	reservationRepository.On("FetchReservationById", "reservation1").Return(&models.ReservationData{ID: "reservation1", UserId: "user1"}, nil)
	notificationRepository.On("CreateNotification", "user1", "reservation1", "New reservation confirmed").Return(errors.New("failed to create notification"))

	// サービス層メソッドの実行
	err := notificationService.CreateNotification("user1", "reservation1", "New reservation confirmed")

	// エラーチェック
	assert.Error(t, err)
	assert.Equal(t, "failed to create notification", err.Error())

	// モックが期待通りに呼び出されてないか確認
	userRepository.AssertCalled(t, "FetchUserById", "user1")
	reservationRepository.AssertCalled(t, "FetchReservationById", "reservation1")
	notificationRepository.AssertCalled(t, "CreateNotification", "user1", "reservation1", "New reservation confirmed")
}
