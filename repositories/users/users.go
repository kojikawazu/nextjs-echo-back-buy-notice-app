package repositories_users

import (
	"backend/models"
	"backend/supabase"
	"log"
)

// Supabaseから全ユーザーを取得し、ユーザーリストを返す。
// 失敗した場合はエラーを返す。
func (r *UserRepositoryImpl) FetchUsers() ([]models.UserData, error) {
	log.Println("Fetching users from Supabase...")

	query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        ORDER BY created_at DESC
    `

	// Supabaseからクエリを実行し、全ユーザーを取得
	rows, err := supabase.Pool.Query(supabase.Ctx, query)
	if err != nil {
		log.Printf("Failed to fetch users: %v", err)
		return nil, err
	}
	log.Println("Fetched users successfully")
	defer rows.Close()

	var users []models.UserData

	// 結果をスキャンしてユーザーデータをリストに追加
	for rows.Next() {
		var user models.UserData
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan user: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		log.Printf("Failed to fetch users: %v", rows.Err())
		return nil, rows.Err()
	}

	log.Printf("Fetched %d users", len(users))
	return users, nil
}

// 指定されたメールアドレスとパスワードでユーザーを取得する。
// ユーザーが見つからない場合、エラーを返す。
func (r *UserRepositoryImpl) FetchUserByEmailAndPassword(email, password string) (*models.UserData, error) {
	log.Printf("Fetching user from Supabase by email: %s\n", email)

	query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        WHERE email = $1 AND password = $2
        LIMIT 1
    `

	// Supabaseからクエリを実行し、条件に一致するユーザーを取得
	row := supabase.Pool.QueryRow(supabase.Ctx, query, email, password)

	// 取得した結果をスキャン
	var user models.UserData
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		log.Printf("User not found or failed to fetch user: %v", err)
		return nil, err
	}

	log.Printf("Fetched user successfully: %v", user)
	return &user, nil
}

// 指定されたIDに対応するユーザーを取得する。
// ユーザーが見つからない場合、エラーを返す。
func (r *UserRepositoryImpl) FetchUserById(id string) (*models.UserData, error) {
	log.Printf("Checking if user exists with id: %s\n", id)

	query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        WHERE id = $1
        LIMIT 1
    `

	// Supabaseからクエリを実行し、条件に一致するユーザーを取得
	row := supabase.Pool.QueryRow(supabase.Ctx, query, id)

	// ユーザーをスキャン
	var user models.UserData
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Printf("User not found or error fetching user: %v", err)
		return nil, err
	}

	log.Printf("User found: %v", user)
	return &user, nil
}

// 指定されたメールアドレスに対応するユーザーを取得する。
// ユーザーが見つからない場合、エラーを返す。
func (r *UserRepositoryImpl) FetchUserByEmail(email string) (*models.UserData, error) {
	log.Printf("Checking if user exists with email: %s\n", email)

	query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        WHERE email = $1
        LIMIT 1
    `

	// Supabaseからクエリを実行し、条件に一致するユーザーを取得
	row := supabase.Pool.QueryRow(supabase.Ctx, query, email)

	// ユーザーをスキャン
	var user models.UserData
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Printf("User not found or error fetching user: %v", err)
		return nil, err
	}

	log.Printf("User found: %v", user)
	return &user, nil
}

// 新しいユーザーをデータベースに追加する。
// 成功した場合はnilを返し、失敗した場合はエラーを返す。
func (r *UserRepositoryImpl) CreateUser(name, email, password string) error {
	log.Printf("Creating new user with email: %s\n", email)

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
        INSERT INTO users (name, email, password, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
    `

	// ユーザーを挿入
	_, err = tx.Exec(supabase.Ctx, query, name, email, password)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return err
	}

	log.Println("User created successfully")
	return nil
}
