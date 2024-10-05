package repositories_reservations

import "backend/models"

// ReservationRepositoryインターフェース
type ReservationRepository interface {
	FetchReservations() ([]models.ReservationData, error)
	FetchReservationById(id string) (*models.ReservationData, error)
	FetchReservationByUserId(userId string) (*models.ReservationData, error)
	CreateReservation(userId, reservationDate string, numPeople int, specialRequest, status string) (string, error)
}

// ReservationRepositoryImplはReservationRepositoryインターフェースを実装する
type ReservationRepositoryImpl struct{}

func NewReservationRepository() ReservationRepository {
	return &ReservationRepositoryImpl{}
}
