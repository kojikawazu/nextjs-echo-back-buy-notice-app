package models

import "time"

// 予約の情報を表すデータ構造
// 各フィールドには、JSONおよびデータベースのタグを指定。
type ReservationData struct {
	ID              string    `json:"id" db:"id"`                             // UUID型
	UserId          string    `json:"user_id" db:"user_id"`                   // ユーザーID
	ReservationDate time.Time `json:"reservation_date" db:"reservation_date"` // 予約日
	NumPeople       int       `json:"num_people" db:"num_people"`             // 予約人数
	SpecialRequest  string    `json:"special_request" db:"special_request"`   // 特別なリクエスト
	Status          string    `json:"status" db:"status"`                     // 予約ステータス
	CreatedAt       time.Time `json:"created_at" db:"created_at"`             // タイムスタンプ
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`             // タイムスタンプ
}
