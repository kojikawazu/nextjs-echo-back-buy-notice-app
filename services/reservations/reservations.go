package services_reservations

import (
	"backend/models"
	"errors"
	"log"
	"time"
)

// Supabaseから全予約情報を取得し、予約情報リストを返す。
// 失敗した場合はエラーを返す。
func (s *ReservationServiceImpl) FetchReservations() ([]models.ReservationData, error) {
	return s.ReservationRepository.FetchReservations()
}

// 指定されたIDに対応する予約情報を取得する。
// 予約情報が見つからない場合、エラーを返す。
func (s *ReservationServiceImpl) FetchReservationById(id string) (*models.ReservationData, error) {
	return s.ReservationRepository.FetchReservationById(id)
}

// 指定されたユーザーIDに対応する予約情報を取得する。
// 予約情報が見つからない場合、エラーを返す。
func (s *ReservationServiceImpl) FetchReservationByUserId(userId string) (*models.ReservationData, error) {
	// バリデーション：userIdが空でないことを確認
	if userId == "" {
		log.Printf("userId is required")
		return nil, errors.New("userId is required")
	}
	log.Println("userId is valid")

	// サービス層からユーザーデータを取得
	reservation, err := s.ReservationRepository.FetchReservationByUserId(userId)
	if err != nil {
		log.Printf("Error fetching reservation: %v", err)
		return nil, errors.New("reservation not found")
	}

	return reservation, nil
}

// 新しい予約情報をデータベースに追加する。
// 成功した場合はnilを返し、失敗した場合はエラーを返す。
func (s *ReservationServiceImpl) CreateReservation(userId, reservationDate string, numPeople int, specialRequest, status string) (string, error) {
	// バリデーション: 必須フィールドが空でないか確認
	if reservationDate == "" || numPeople <= 0 {
		log.Printf("UserID, reservation date, and num_people are required")
		return "", errors.New("userID, reservation date, and num_people are required")
	}

	// 予約日が正しいフォーマットか確認
	_, err := time.Parse("2006-01-02 15:04:05", reservationDate)
	if err != nil {
		log.Printf("Invalid reservation date format: %v", err)
		return "", errors.New("invalid reservation date format. Use 'YYYY-MM-DD HH:MM:SS'")
	}

	// ユーザーが存在するか確認
	existingUser, err := s.UserRepository.FetchUserById(userId)
	if err != nil || existingUser == nil {
		log.Printf("User not found: %s", userId)
		return "", errors.New("user not found")
	}

	// ステータスが指定されていない場合、デフォルトで"pending"とする
	if status == "" {
		status = "pending"
	}

	log.Println("Request body is valid")

	// 予約を作成する
	reservationId, err := s.ReservationRepository.CreateReservation(userId, reservationDate, numPeople, specialRequest, status)
	if err != nil {
		log.Printf("Error creating reservation: %v", err)
		return "", errors.New("failed to create reservation")
	}

	return reservationId, nil
}
