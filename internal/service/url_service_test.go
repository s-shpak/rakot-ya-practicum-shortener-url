package service

import (
	"errors"
	"github.com/rakot9/ya-practicum-shortener-url/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
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

func TestURLService_ShortenURL_Success(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	targetURL := "https://yandex.ru"

	// Настраиваем мок:
	// 1. Метод Find вернет ErrNotFound, значит ключ свободен.
	// 2. Метод Save успешно сохранит пару ключ-значение.
	mockRepo.On("Find", mock.AnythingOfType("string")).Return("", repository.ErrNotFound).Once()
	mockRepo.On("Save", mock.AnythingOfType("string"), targetURL).Return(nil).Once()

	key, err := service.ShortenURL(targetURL)

	assert.NoError(t, err)
	assert.NotEmpty(t, key)
	// Ключ должен быть длиной от 6 до 9 символов (исходя из логики generateKey)
	assert.GreaterOrEqual(t, len(key), 6)
	assert.LessOrEqual(t, len(key), 9)
	mockRepo.AssertExpectations(t)
}

func TestURLService_ShortenURL_CollisionResolved(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	targetURL := "https://yandex.ru"

	// Симулируем коллизию:
	// 1-я попытка: ключ найден в БД (коллизия)
	mockRepo.On("Find", mock.AnythingOfType("string")).Return("https://old-url.com", nil).Once()

	// 2-я попытка: новый сгенерированный ключ свободен
	mockRepo.On("Find", mock.AnythingOfType("string")).Return("", repository.ErrNotFound).Once()
	mockRepo.On("Save", mock.AnythingOfType("string"), targetURL).Return(nil).Once()

	key, err := service.ShortenURL(targetURL)

	assert.NoError(t, err)
	assert.NotEmpty(t, key)
	mockRepo.AssertExpectations(t)
}

func TestURLService_ShortenURL_AttemptsExceeded(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	targetURL := "https://yandex.ru"

	// Симулируем жесткую коллизию: все 10 раз Find возвращает существующий URL
	mockRepo.On("Find", mock.AnythingOfType("string")).Return("https://existing.com", nil).Times(10)

	key, err := service.ShortenURL(targetURL)

	assert.Empty(t, key)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrCountAttemptExceed)
	mockRepo.AssertExpectations(t)
}

func TestURLService_ShortenURL_RepoFindError(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	targetURL := "https://yandex.ru"
	dbErr := errors.New("connection timeout")

	// Симулируем критическую ошибку БД при поиске ключа
	mockRepo.On("Find", mock.AnythingOfType("string")).Return("", dbErr).Once()

	key, err := service.ShortenURL(targetURL)

	assert.Empty(t, key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check key existence")
	mockRepo.AssertExpectations(t)
}

func TestURLService_ShortenURL_RepoSaveError(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	targetURL := "https://yandex.ru"
	dbErr := errors.New("write permission denied")

	// Ключ свободен, но при сохранении падает ошибка
	mockRepo.On("Find", mock.AnythingOfType("string")).Return("", repository.ErrNotFound).Once()
	mockRepo.On("Save", mock.AnythingOfType("string"), targetURL).Return(dbErr).Once()

	key, err := service.ShortenURL(targetURL)

	assert.Empty(t, key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save url")
	mockRepo.AssertExpectations(t)
}

func TestURLService_GetURL_Success(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	expectedKey := "aB3dE6"
	expectedURL := "https://yandex.ru"

	mockRepo.On("Find", expectedKey).Return(expectedURL, nil).Once()

	url, err := service.GetURL(expectedKey)

	assert.NoError(t, err)
	assert.Equal(t, expectedURL, url)
	mockRepo.AssertExpectations(t)
}

func TestURLService_GetURL_NotFound(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	expectedKey := "nonexistent"

	mockRepo.On("Find", expectedKey).Return("", repository.ErrNotFound).Once()

	url, err := service.GetURL(expectedKey)

	assert.Empty(t, url)
	assert.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrNotFound)
	mockRepo.AssertExpectations(t)
}
