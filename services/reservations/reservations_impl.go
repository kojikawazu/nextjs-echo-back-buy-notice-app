package services_reservations

import (
	"backend/models"
	repositories_reservations "backend/repositories/reservations"
	repositories_users "backend/repositories/users"
)

// ReservationServiceインターフェース
type ReservationService interface {
	FetchReservations() ([]models.ReservationData, error)
	FetchReservationById(id string) (*models.ReservationData, error)
	FetchReservationByUserId(userId string) (*models.ReservationData, error)
	CreateReservation(userId, reservationDate string, numPeople int, specialRequest, status string) (string, error)
}

// ReservationServiceImplはReservationServiceインターフェースを実装する
type ReservationServiceImpl struct {
	UserRepository        repositories_users.UserRepository
	ReservationRepository repositories_reservations.ReservationRepository
}

func NewReservationService(
	userRepository repositories_users.UserRepository,
	reservationRepository repositories_reservations.ReservationRepository,
) ReservationService {
	return &ReservationServiceImpl{
		UserRepository:        userRepository,
		ReservationRepository: reservationRepository,
	}
}
