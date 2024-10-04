package repositories_notifications

import "backend/models"

// NotificationRepositoryインターフェース
type NotificationRepository interface {
	FetchNotifications() ([]models.NotificationData, error)
	CreateNotification(userId, reservationId, message string) error
}

// NotificationRepositoryImplはNotificationRepositoryインターフェースを実装する
type NotificationRepositoryImpl struct{}

func NewNotificationRepository() NotificationRepository {
	return &NotificationRepositoryImpl{}
}
