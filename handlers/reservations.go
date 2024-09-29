package handlers

import (
	"backend/services"
	"backend/websocket"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// 全予約情報を取得し、JSON形式で返すハンドラー
// 予約情報取得に失敗した場合、500エラーを返す。
func GetReservations(c echo.Context) error {
	log.Println("Fetching reservations...")

	// サービス層で予約情報一覧を取得
	reservations, err := services.FetchReservations()
	if err != nil {
		log.Printf("Error fetching reservations from Supabase: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch reservations",
		})
	}

	log.Println("Fetched reservations successfully")
	return c.JSON(http.StatusOK, reservations)
}

// パスパラメータで指定されたユーザーIDで予約情報を取得する。
// データベースに該当予約情報がいない場合、404エラーを返す。
func GetReservationByUserId(c echo.Context) error {
	log.Println("Fetching reservation by userId...")

	// パスパラメータからuserIdを取得
	userId := c.Param("user_id")

	// バリデーション：userIdが空でないことを確認
	if userId == "" {
		log.Printf("userId is required")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "userId is required",
		})
	}
	log.Println("userId is valid")

	// サービス層からユーザーデータを取得
	reservation, err := services.FetchReservationByUserId(userId)
	if err != nil {
		log.Printf("Error fetching reservation: %v", err)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Reservation not found",
		})
	}

	log.Println("Fetched reservation successfully")
	return c.JSON(http.StatusOK, reservation)
}

// 新しい予約情報を追加するハンドラー
func AddReservation(c echo.Context) error {
	log.Println("Creating new reservation...")

	// リクエストボディからデータを取得
	type RequestBody struct {
		UserID          string `json:"user_id"`          // ユーザーID
		ReservationDate string `json:"reservation_date"` // 予約日
		NumPeople       int    `json:"num_people"`       // 人数
		SpecialRequest  string `json:"special_request"`  // 特別リクエスト
		Status          string `json:"status"`           // ステータス
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
	if reqBody.UserID == "" || reqBody.ReservationDate == "" || reqBody.NumPeople <= 0 {
		log.Printf("UserID, reservation date, and num_people are required")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "UserID, reservation date, and num_people are required",
		})
	}

	// 予約日が正しいフォーマットか確認
	_, err := time.Parse("2006-01-02 15:04:05", reqBody.ReservationDate)
	if err != nil {
		log.Printf("Invalid reservation date format: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid reservation date format. Use 'YYYY-MM-DD HH:MM:SS'",
		})
	}

	// ユーザーが存在するか確認
	existingUser, err := services.FetchUserById(reqBody.UserID)
	if err != nil || existingUser == nil {
		log.Printf("User not found: %s", reqBody.UserID)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	// ステータスが指定されていない場合、デフォルトで"pending"とする
	if reqBody.Status == "" {
		reqBody.Status = "pending"
	}

	log.Println("Request body is valid")

	// 予約を作成する
	reservationId, err := services.CreateReservation(reqBody.UserID, reqBody.ReservationDate, reqBody.NumPeople, reqBody.SpecialRequest, reqBody.Status)
	if err != nil {
		log.Printf("Error creating reservation: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create reservation",
		})
	}
	log.Println("Reservation created successfully")

	// 予約が成功したので通知情報を作成
	notificationMessage := "New reservation created for user " + reqBody.UserID
	err = services.CreateNotification(reqBody.UserID, reservationId, notificationMessage)
	if err != nil {
		log.Printf("Error creating notification: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create notification",
		})
	}
	log.Println("Notification created successfully")

	// Redisに通知メッセージをパブリッシュ
	err = websocket.PublishToRedis("reservation-notifications", notificationMessage)
	if err != nil {
		log.Printf("Failed to publish notification to Redis: %v", err)
	}
	log.Println("Published reservation notification to Redis successfully")

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Reservation created successfully",
	})
}
