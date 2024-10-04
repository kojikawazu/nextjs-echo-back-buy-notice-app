package services_notifications

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

// MockNotificationService is the mock implementation for NotificationService
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) FetchNotifications() ([]models.NotificationData, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]models.NotificationData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNotificationService) CreateNotification(userID, reservationID, message string) error {
	args := m.Called(userID, reservationID, message)
	return args.Error(0)
}
