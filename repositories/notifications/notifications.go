package repositories_notifications

import (
	"backend/models"
	"backend/supabase"
	"log"
	"time"
)

// Supabaseから全通知情報を取得し、通知情報リストを返す。
// 失敗した場合はエラーを返す。
func (r *NotificationRepositoryImpl) FetchNotifications() ([]models.NotificationData, error) {
	log.Println("Fetching notifications from Supabase...")

	query := `
        SELECT id, user_id, reservation_id, message, created_at
        FROM notifications
        ORDER BY created_at DESC
    `

	// Supabaseからクエリを実行し、全通知情報を取得
	rows, err := supabase.Pool.Query(supabase.Ctx, query)
	if err != nil {
		log.Printf("Failed to fetch notifications: %v", err)
		return nil, err
	}
	log.Println("Fetched notifications successfully")
	defer rows.Close()

	var notifications []models.NotificationData

	// 結果をスキャンしてユーザーデータをリストに追加
	for rows.Next() {
		var notification models.NotificationData
		err := rows.Scan(
			&notification.ID,
			&notification.UserId,
			&notification.ReservationId,
			&notification.Message,
			&notification.CreatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan notification: %v", err)
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	if rows.Err() != nil {
		log.Printf("Failed to fetch reservations: %v", rows.Err())
		return nil, rows.Err()
	}

	log.Printf("Fetched %d notifications", len(notifications))
	return notifications, nil
}

// 新しい通知をデータベースに追加する。
// 成功した場合はnilを返し、失敗した場合はエラーを返す。
func (r *NotificationRepositoryImpl) CreateNotification(userId, reservationId, message string) error {
	log.Printf("Creating new notification for userId: %s\n", userId)

	// トランザクションの開始
	tx, err := supabase.Pool.Begin(supabase.Ctx)
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return err
	}

	// トランザクションが成功または失敗した場合にコミットまたはロールバックを行う
	defer func() {
		if err != nil {
			log.Println("Rolling back transaction...")
			if rollbackErr := tx.Rollback(supabase.Ctx); rollbackErr != nil {
				log.Printf("Failed to rollback transaction: %v", rollbackErr)
			}
			return
		}

		log.Println("Committing transaction...")
		if commitErr := tx.Commit(supabase.Ctx); commitErr != nil {
			log.Printf("Failed to commit transaction: %v", commitErr)
		}
	}()

	// 通知を挿入するSQLクエリ
	query := `
        INSERT INTO notifications (user_id, reservation_id, message, created_at)
        VALUES ($1, $2, $3, $4)
    `

	// 現在の時刻を取得
	createdAt := time.Now()

	// 通知をデータベースに挿入
	_, err = tx.Exec(supabase.Ctx, query, userId, reservationId, message, createdAt)
	if err != nil {
		log.Printf("Failed to create notification: %v", err)
		return err
	}

	log.Println("Notification created successfully")
	return nil
}
