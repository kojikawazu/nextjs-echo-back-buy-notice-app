package services_users

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

// MockUserService is the mock implementation for UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) FetchUsers() ([]models.UserData, error) {
	args := m.Called()
	return args.Get(0).([]models.UserData), args.Error(1)
}

func (m *MockUserService) FetchUserByEmailAndPassword(email, password string) (*models.UserData, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserData), args.Error(1)
}

func (m *MockUserService) FetchUserByEmail(email string) (*models.UserData, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserData), args.Error(1)
}

func (m *MockUserService) FetchUserById(id string) (*models.UserData, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserData), args.Error(1)
}

func (m *MockUserService) CreateUser(name, email, password string) error {
	args := m.Called(name, email, password)
	return args.Error(0)
}
