package services_notifications

import (
	"backend/models"
	repositories_notifications "backend/repositories/notifications"
	services_reservations "backend/services/reservations"
	services_users "backend/services/users"
)

// NotificationServiceインターフェース
type NotificationService interface {
	FetchNotifications() ([]models.NotificationData, error)
	CreateNotification(userId, reservationId, message string) error
}

// NotificationServiceImplはNotificationServiceインターフェースを実装する
type NotificationServiceImpl struct {
	UserService            services_users.UserService
	ReservationService     services_reservations.ReservationService
	NotificationRepository repositories_notifications.NotificationRepository
}

func NewNotificationService(
	userService services_users.UserService,
	reservationService services_reservations.ReservationService,
	notificationRepository repositories_notifications.NotificationRepository,
) NotificationService {
	return &NotificationServiceImpl{
		UserService:            userService,
		ReservationService:     reservationService,
		NotificationRepository: notificationRepository,
	}
}
