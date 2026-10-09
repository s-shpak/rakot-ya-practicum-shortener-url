package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/rakot9/ya-practicum-shortener-url/internal/repository"
	"math/big"
)

const base62Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var ErrCountAttemptExceed = errors.New("count attempt exceed")

// URLService управляет бизнес-логикой генерации ключей и координирует работу с репозиторием.
type URLService struct {
	repo repository.URLRepository
}

// NewURLService принимает любой репозиторий, удовлетворяющий интерфейсу URLRepository.
func NewURLService(repo repository.URLRepository) *URLService {
	return &URLService{
		repo: repo,
	}
}

// ShortenURL генерирует уникальный ключ и сохраняет его в репозиторий.
func (s *URLService) ShortenURL(url string) (string, error) {
	for i := 10; i > 0; i-- {
		key, err := s.generateKey()
		if err != nil {
			return "", fmt.Errorf("failed to generate random id: %w", err)
		}

		// Проверяем, существует ли уже ключ в репозитории
		_, err = s.repo.Find(key)
		if errors.Is(err, repository.ErrNotFound) {
			// Ключ свободен, сохраняем и выходим из цикла
			if err := s.repo.Save(key, url); err != nil {
				return "", fmt.Errorf("failed to save url: %w", err)
			}
			return key, nil
		}
		if err != nil {
			// Если произошла другая ошибка (например, упала БД), возвращаем её
			return "", fmt.Errorf("failed to check key existence: %w", err)
		}
	}

	return "", fmt.Errorf("failed generate key 10 attempt exceed for url %s error %w", url, ErrCountAttemptExceed)
}

// GetURL возвращает оригинальный URL по его короткому ключу.
func (s *URLService) GetURL(key string) (string, error) {
	return s.repo.Find(key)
}

// Вспомогательный метод генерации ключа остался внутри сервиса, так как это бизнес-логика.
func (s *URLService) generateKey() (string, error) {
	lengthBig, err := rand.Int(rand.Reader, big.NewInt(4)) // случайное число от 0 до 3
	if err != nil {
		return "", err
	}
	length := int(lengthBig.Int64()) + 6

	b := make([]byte, length)
	alphabetLen := big.NewInt(int64(len(base62Alphabet)))

	for i := 0; i < length; i++ {
		idx, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", err
		}
		b[i] = base62Alphabet[idx.Int64()]
	}

	return string(b), nil
}
