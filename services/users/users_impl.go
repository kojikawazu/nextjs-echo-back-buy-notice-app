package services_users

import (
	"backend/models"
	repositories_users "backend/repositories/users"
)

// UserServiceインターフェース
type UserService interface {
	FetchUsers() ([]models.UserData, error)
	FetchUserByEmailAndPassword(email, password string) (*models.UserData, error)
	FetchUserById(id string) (*models.UserData, error)
	FetchUserByEmail(email string) (*models.UserData, error)
	CreateUser(name, email, password string) error
}

// UserServiceImplはUserServiceインターフェースを実装する
type UserServiceImpl struct {
	UserRepository repositories_users.UserRepository
}

func NewUserService(
	userRepository repositories_users.UserRepository,
) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
	}
}
