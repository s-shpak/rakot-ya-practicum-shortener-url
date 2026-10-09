package repository

import (
	"github.com/stretchr/testify/mock"
)

type MockURLRepository struct {
	mock.Mock
}

func (m *MockURLRepository) Find(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockURLRepository) Save(key, url string) error {
	args := m.Called(key, url)
	return args.Error(0)
}
