package repositories_users

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FetchUsers() ([]models.UserData, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]models.UserData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FetchUserById(id string) (*models.UserData, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.UserData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FetchUserByEmailAndPassword(email, password string) (*models.UserData, error) {
	args := m.Called(email, password)
	if args.Get(0) != nil {
		return args.Get(0).(*models.UserData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FetchUserByEmail(email string) (*models.UserData, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*models.UserData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) CreateUser(name, email, password string) error {
	args := m.Called(name, email, password)
	return args.Error(0)
}
