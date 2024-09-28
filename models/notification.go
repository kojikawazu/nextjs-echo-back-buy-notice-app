package models

// 通知の情報を表すデータ構造
// 各フィールドには、JSONおよびデータベースのタグを指定。
type NotificationData struct {
	ID            string `json:"id" db:"id"`                         // UUID型
	UserId        string `json:"user_id" db:"user_id"`               // ユーザーID
	ReservationId string `json:"reservation_id" db:"reservation_id"` // 予約ID
	Message       string `json:"message" db:"message"`               // 通知メッセージ
	CreatedAt     string `json:"created_at" db:"created_at"`         // タイムスタンプ
}
