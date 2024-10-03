package services_users

import "backend/models"

// UserServiceインターフェース
type UserService interface {
	FetchUsers() ([]models.UserData, error)
	FetchUserByEmailAndPassword(email, password string) (*models.UserData, error)
	FetchUserById(id string) (*models.UserData, error)
	FetchUserByEmail(email string) (*models.UserData, error)
	CreateUser(name, email, password string) error
}

// UserServiceImplはUserServiceインターフェースを実装する
type UserServiceImpl struct{}

func (u *UserServiceImpl) FetchUsers() ([]models.UserData, error) {
	return FetchUsers()
}

func (u *UserServiceImpl) FetchUserByEmailAndPassword(email, password string) (*models.UserData, error) {
	return FetchUserByEmailAndPassword(email, password)
}

func (u *UserServiceImpl) FetchUserById(id string) (*models.UserData, error) {
	return FetchUserById(id)
}

func (u *UserServiceImpl) FetchUserByEmail(email string) (*models.UserData, error) {
	return FetchUserByEmail(email)
}

func (u *UserServiceImpl) CreateUser(name, email, password string) error {
	return CreateUser(name, email, password)
}
