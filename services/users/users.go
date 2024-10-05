package services_users

import (
	"backend/models"
	"database/sql"
	"errors"
	"log"
	"net/mail"
)

// Supabaseから全ユーザーを取得し、ユーザーリストを返す。
// 失敗した場合はエラーを返す。
func (s *UserServiceImpl) FetchUsers() ([]models.UserData, error) {
	return s.UserRepository.FetchUsers()
}

// 指定されたメールアドレスとパスワードでユーザーを取得する。
// ユーザーが見つからない場合、エラーを返す。
func (s *UserServiceImpl) FetchUserByEmailAndPassword(email, password string) (*models.UserData, error) {
	// バリデーション：emailとpasswordが空でないことを確認
	if email == "" || password == "" {
		log.Printf("Email and password are required")
		return nil, errors.New("email and password are required")
	}
	// バリデーション：emailが有効な形式であることを確認
	if _, err := mail.ParseAddress(email); err != nil {
		log.Printf("Invalid email format: %v", err)
		return nil, errors.New("invalid email format")
	}
	log.Println("Email and password are valid")

	user, err := s.UserRepository.FetchUserByEmailAndPassword(email, password)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User not found for email: %s", email)
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// 指定されたIDに対応するユーザーを取得する。
// ユーザーが見つからない場合、エラーを返す。
func (s *UserServiceImpl) FetchUserById(id string) (*models.UserData, error) {
	return s.UserRepository.FetchUserById(id)
}

// 指定されたメールアドレスに対応するユーザーを取得する。
// ユーザーが見つからない場合、エラーを返す。
func (s *UserServiceImpl) FetchUserByEmail(email string) (*models.UserData, error) {
	return s.UserRepository.FetchUserByEmail(email)
}

// 新しいユーザーをデータベースに追加する。
// 成功した場合はnilを返し、失敗した場合はエラーを返す。
func (s *UserServiceImpl) CreateUser(name, email, password string) error {
	// バリデーション: 名前、Email、パスワードが空でないかを確認
	if name == "" || email == "" || password == "" {
		log.Printf("Name, email and password are required")
		return errors.New("name, email and password are required")
	}

	// Eメール形式のバリデーション
	if _, err := mail.ParseAddress(email); err != nil {
		log.Printf("Invalid email format: %v", err)
		return errors.New("invalid email format")
	}
	log.Println("Name, email and password are valid")

	// 既存ユーザーの確認
	existingUser, err := s.UserRepository.FetchUserByEmail(email)
	if err == nil && existingUser != nil {
		// ユーザーが既に存在する場合はスキップ
		log.Printf("User already exists: %v", existingUser)
		return errors.New("user already exists")
	}
	log.Println("User does not exist")

	err = s.UserRepository.CreateUser(name, email, password)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return errors.New("failed to create user")
	}

	return nil
}
