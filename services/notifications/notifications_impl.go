package services_notifications

import "backend/models"

// NotificationServiceインターフェース
type NotificationService interface {
	FetchNotifications() ([]models.NotificationData, error)
	CreateNotification(userId, reservationId, message string) error
}

// NotificationServiceImplはNotificationServiceインターフェースを実装する
type NotificationServiceImpl struct{}

func (u *NotificationServiceImpl) FetchNotifications() ([]models.NotificationData, error) {
	return FetchNotifications()
}

func (u *NotificationServiceImpl) CreateNotification(userId, reservationId, message string) error {
	return CreateNotification(userId, reservationId, message)
}
