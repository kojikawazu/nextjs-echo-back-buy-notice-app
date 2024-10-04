package repositories_notifications

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

// MockNotificationRepository is a mock implementation of NotificationRepository
type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) FetchNotifications() ([]models.NotificationData, error) {
	args := m.Called()
	return args.Get(0).([]models.NotificationData), args.Error(1)
}

func (m *MockNotificationRepository) CreateNotification(userId, reservationId, message string) error {
	args := m.Called(userId, reservationId, message)
	return args.Error(0)
}
