package handlers_notifications

import (
	services_notifications "backend/services/notifications"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type NotificationHandler struct {
	NotificationService services_notifications.NotificationService
}

// コンストラクタ
func NewNotificationHandler(notificationService services_notifications.NotificationService) *NotificationHandler {
	return &NotificationHandler{
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

	// 通知を作成
	err := h.NotificationService.CreateNotification(reqBody.UserID, reqBody.ReservationID, reqBody.Message)
	if err != nil {
		switch err.Error() {
		case "userID, ReservationID, and message are required":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "UserID, ReservationID, and message are required",
			})
		case "user not found":
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		case "reservation not found":
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Reservation not found",
			})
		case "failed to create notification":
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to create notification",
			})
		default:
			log.Printf("failed to create reservation: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to create reservation",
			})
		}
	}

	log.Println("Notification created successfully")
	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Notification created successfully",
	})
}
