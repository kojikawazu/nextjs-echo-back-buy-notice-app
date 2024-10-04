package services_reservations

import "backend/models"

// ReservationServiceインターフェース
type ReservationService interface {
	FetchReservations() ([]models.ReservationData, error)
	FetchReservationById(id string) (*models.ReservationData, error)
	FetchReservationByUserId(userId string) (*models.ReservationData, error)
	CreateReservation(userId, reservationDate string, numPeople int, specialRequest, status string) (string, error)
}

// ReservationServiceImplはReservationServiceインターフェースを実装する
type ReservationServiceImpl struct{}

func (u *ReservationServiceImpl) FetchReservations() ([]models.ReservationData, error) {
	return FetchReservations()
}

func (u *ReservationServiceImpl) FetchReservationById(id string) (*models.ReservationData, error) {
	return FetchReservationById(id)
}

func (u *ReservationServiceImpl) FetchReservationByUserId(userId string) (*models.ReservationData, error) {
	return FetchReservationByUserId(userId)
}

func (u *ReservationServiceImpl) CreateReservation(userId, reservationDate string, numPeople int, specialRequest, status string) (string, error) {
	return CreateReservation(userId, reservationDate, numPeople, specialRequest, status)
}
