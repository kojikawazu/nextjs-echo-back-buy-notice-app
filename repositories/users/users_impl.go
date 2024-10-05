package repositories_users

import "backend/models"

// UserRepositoryインターフェース
type UserRepository interface {
	FetchUsers() ([]models.UserData, error)
	FetchUserByEmailAndPassword(email, password string) (*models.UserData, error)
	FetchUserById(id string) (*models.UserData, error)
	FetchUserByEmail(email string) (*models.UserData, error)
	CreateUser(name, email, password string) error
}

// UserRepositoryImplはUserRepositoryインターフェースを実装する
type UserRepositoryImpl struct{}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}
