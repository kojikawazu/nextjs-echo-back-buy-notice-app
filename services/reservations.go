package services

import (
	"backend/models"
	"backend/supabase"
	"log"
)

// Supabaseから全予約情報を取得し、予約情報リストを返す。
// 失敗した場合はエラーを返す。
func FetchReservations() ([]models.ReservationData, error) {
	log.Println("Fetching reservations from Supabase...")

	query := `
        SELECT id, user_id, reservation_date, num_people, special_request, status, created_at, updated_at
        FROM reservations
        ORDER BY created_at DESC
    `

	rows, err := supabase.Pool.Query(supabase.Ctx, query)
	if err != nil {
		log.Printf("Failed to fetch reservations: %v", err)
		return nil, err
	}
	log.Println("Fetched reservations successfully")
	defer rows.Close()

	var reservations []models.ReservationData

	// 結果をスキャンしてユーザーデータをリストに追加
	for rows.Next() {
		var reservation models.ReservationData
		err := rows.Scan(
			&reservation.ID,
			&reservation.UserId,
			&reservation.ReservationDate,
			&reservation.NumPeople,
			&reservation.SpecialRequest,
			&reservation.Status,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan reservation: %v", err)
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	if rows.Err() != nil {
		log.Printf("Failed to fetch reservations: %v", rows.Err())
		return nil, rows.Err()
	}

	log.Printf("Fetched %d reservations", len(reservations))
	return reservations, nil
}

// 指定されたユーザーIDに対応する予約情報を取得する。
// 予約情報が見つからない場合、エラーを返す。
func FetchReservationByUserId(userId string) (*models.ReservationData, error) {
	log.Printf("Checking if reservation exists with userId: %s\n", userId)

	query := `
        SELECT id, user_id, reservation_date, num_people, special_request, status, created_at, updated_at
        FROM reservations
        WHERE user_id = $1
    `

	row := supabase.Pool.QueryRow(supabase.Ctx, query, userId)

	var reservation models.ReservationData
	err := row.Scan(&reservation.ID, &reservation.UserId, &reservation.ReservationDate, &reservation.NumPeople, &reservation.SpecialRequest, &reservation.Status, &reservation.CreatedAt, &reservation.UpdatedAt)
	if err != nil {
		log.Printf("Reservation not found or error fetching reservation: %v", err)
		return nil, err
	}

	log.Printf("Reservation found: %v", reservation)
	return &reservation, nil
}

// 新しい予約情報をデータベースに追加する。
// 成功した場合はnilを返し、失敗した場合はエラーを返す。
func CreateReservation(userId, reservationDate string, numPeople int, specialRequest, status string) error {
	log.Printf("Creating new reservation for userId: %s\n", userId)

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

	query := `
        INSERT INTO reservations (user_id, reservation_date, num_people, special_request, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
    `

	_, err = tx.Exec(supabase.Ctx, query, userId, reservationDate, numPeople, specialRequest, status)
	if err != nil {
		log.Printf("Failed to create reservation: %v", err)
		return err
	}

	log.Println("Reservation created successfully")

	return nil
}
