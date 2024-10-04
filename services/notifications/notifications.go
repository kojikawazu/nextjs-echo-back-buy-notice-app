package services_notifications

import (
	"backend/models"
	"errors"
	"log"
)

// Supabaseから全通知情報を取得し、通知情報リストを返す。
// 失敗した場合はエラーを返す。
func (s *NotificationServiceImpl) FetchNotifications() ([]models.NotificationData, error) {
	return s.NotificationRepository.FetchNotifications()
}

// 新しい通知をデータベースに追加する。
// 成功した場合はnilを返し、失敗した場合はエラーを返す。
func (s *NotificationServiceImpl) CreateNotification(userId, reservationId, message string) error {
	// バリデーション: 必須フィールドが空でないか確認
	if userId == "" || reservationId == "" || message == "" {
		log.Printf("userID, ReservationID, and message are required")
		return errors.New("userID, ReservationID, and message are required")
	}

	// ユーザーの存在確認
	existingUser, err := s.UserService.FetchUserById(userId)
	if err != nil || existingUser == nil {
		log.Printf("user not found: %s", userId)
		return errors.New("user not found")
	}

	// 予約の存在確認（オプションだが、予約が実在するか確認したい場合）
	existingReservation, err := s.ReservationService.FetchReservationById(reservationId)
	if err != nil || existingReservation == nil {
		log.Printf("reservation not found: %s", reservationId)
		return errors.New("reservation not found")
	}

	return s.NotificationRepository.CreateNotification(userId, reservationId, message)
}
