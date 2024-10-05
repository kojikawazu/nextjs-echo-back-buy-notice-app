package services_reservations

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

// MockReservationService is the mock implementation for ReservationService
type MockReservationService struct {
	mock.Mock
}

func (m *MockReservationService) FetchReservations() ([]models.ReservationData, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]models.ReservationData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockReservationService) FetchReservationById(id string) (*models.ReservationData, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ReservationData), args.Error(1)
}

func (m *MockReservationService) FetchReservationByUserId(userId string) (*models.ReservationData, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ReservationData), args.Error(1)
}

func (m *MockReservationService) CreateReservation(userId, reservationDate string, numPeople int, specialRequest, status string) (string, error) {
	args := m.Called(userId, reservationDate, numPeople, specialRequest, status)
	return args.String(0), args.Error(1)
}
