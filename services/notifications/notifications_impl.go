package services_notifications

import (
	"backend/models"
	repositories_notifications "backend/repositories/notifications"
	repositories_reservations "backend/repositories/reservations"
	repositories_users "backend/repositories/users"
)

// NotificationServiceインターフェース
type NotificationService interface {
	FetchNotifications() ([]models.NotificationData, error)
	CreateNotification(userId, reservationId, message string) error
}

// NotificationServiceImplはNotificationServiceインターフェースを実装する
type NotificationServiceImpl struct {
	UserRepository         repositories_users.UserRepository
	ReservationRepository  repositories_reservations.ReservationRepository
	NotificationRepository repositories_notifications.NotificationRepository
}

func NewNotificationService(
	userRepository repositories_users.UserRepository,
	reservationRepository repositories_reservations.ReservationRepository,
	notificationRepository repositories_notifications.NotificationRepository,
) NotificationService {
	return &NotificationServiceImpl{
		UserRepository:         userRepository,
		ReservationRepository:  reservationRepository,
		NotificationRepository: notificationRepository,
	}
}
