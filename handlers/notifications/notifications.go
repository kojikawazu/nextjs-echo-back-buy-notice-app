package handlers_notifications

import (
	services_notifications "backend/services/notifications"
	services_reservations "backend/services/reservations"
	services_users "backend/services/users"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type NotificationHandler struct {
	UserService         services_users.UserService
	ReservationService  services_reservations.ReservationService
	NotificationService services_notifications.NotificationService
}

// コンストラクタ
func NewNotificationHandler(userService services_users.UserService, reservationService services_reservations.ReservationService, notificationService services_notifications.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		UserService:         userService,
		ReservationService:  reservationService,
		NotificationService: notificationService,
	}
}

// 全通知情報を取得し、JSON形式で返すハンドラー
// 通知情報取得に失敗した場合、500エラーを返す。
func (h *NotificationHandler) GetNotifications(c echo.Context) error {
	log.Println("Fetching notifications...")

	// サービス層で予約情報一覧を取得
	notifications, err := h.NotificationService.FetchNotifications()
	if err != nil {
		log.Printf("Error fetching notifications from Supabase: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch notifications",
		})
	}

	log.Println("Fetched notifications successfully")
	return c.JSON(http.StatusOK, notifications)
}

// 新しい通知を追加するハンドラー
func (h *NotificationHandler) AddNotification(c echo.Context) error {
	log.Println("Creating new notification...")

	// リクエストボディからデータを取得
	type RequestBody struct {
		UserID        string `json:"user_id"`        // ユーザーID
		ReservationID string `json:"reservation_id"` // 予約ID
		Message       string `json:"message"`        // 通知メッセージ
	}

	// リクエストボディをバインド
	var reqBody RequestBody
	if err := c.Bind(&reqBody); err != nil {
		log.Printf("Failed to bind request body: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// バリデーション: 必須フィールドが空でないか確認
	if reqBody.UserID == "" || reqBody.ReservationID == "" || reqBody.Message == "" {
		log.Printf("UserID, ReservationID, and message are required")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "UserID, ReservationID, and message are required",
		})
	}

	// ユーザーの存在確認
	existingUser, err := h.UserService.FetchUserById(reqBody.UserID)
	if err != nil || existingUser == nil {
		log.Printf("User not found: %s", reqBody.UserID)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	// 予約の存在確認（オプションだが、予約が実在するか確認したい場合）
	existingReservation, err := h.ReservationService.FetchReservationById(reqBody.ReservationID)
	if err != nil || existingReservation == nil {
		log.Printf("Reservation not found: %s", reqBody.ReservationID)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Reservation not found",
		})
	}

	// 通知を作成
	err = h.NotificationService.CreateNotification(reqBody.UserID, reqBody.ReservationID, reqBody.Message)
	if err != nil {
		log.Printf("Error creating notification: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create notification",
		})
	}

	log.Println("Notification created successfully")
	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Notification created successfully",
	})
}
