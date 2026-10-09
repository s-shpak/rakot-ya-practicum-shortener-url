package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/rakot9/ya-practicum-shortener-url/internal/repository"
	"github.com/rakot9/ya-practicum-shortener-url/internal/service"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockRepo struct {
	storage map[string]string
}

func (m *mockRepo) Save(key, value string) error {
	m.storage[key] = value
	return nil
}

func (m *mockRepo) Find(key string) (string, error) {
	val, ok := m.storage[key]
	if !ok {
		return "", repository.ErrNotFound
	}
	return val, nil
}

// Инициализация хендлера с mock сервисом
func setupTestHandler() (*Handler, *mockRepo) {
	repo := &mockRepo{storage: make(map[string]string)}

	urlService := service.NewURLService(repo)

	flagEnv := &FlagEnv{
		FlagRunAddr:           "localhost:8080",
		FlagRunShorternerAddr: "http://localhost:8080",
	}

	return NewHandler(urlService, flagEnv), repo
}

func TestHandler_MainPage(t *testing.T) {
	h, _ := setupTestHandler()

	tests := []struct {
		name         string
		contentType  string
		body         string
		wantStatus   int
		wantBodyPart string
	}{
		{
			name:         "Успешное сокращение URL",
			contentType:  "text/plain",
			body:         "https://yandex.ru",
			wantStatus:   http.StatusCreated,
			wantBodyPart: "http://localhost:8080/",
		},
		{
			name:         "Неверный Content-Type",
			contentType:  "application/json",
			body:         "https://yandex.ru",
			wantStatus:   http.StatusBadRequest,
			wantBodyPart: "Неверный http метод запроса",
		},
		{
			name:         "Пустой тип контента",
			contentType:  "",
			body:         "https://yandex.ru",
			wantStatus:   http.StatusBadRequest,
			wantBodyPart: "Неверный http метод запроса",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			if tt.contentType != "" {
				request.Header.Set("Content-Type", tt.contentType)
			}

			recorder := httptest.NewRecorder()
			h.MainPage(recorder, request)

			res := recorder.Result()
			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantStatus, res.StatusCode)
			assert.Contains(t, string(resBody), tt.wantBodyPart)

			if tt.wantStatus == http.StatusCreated {
				assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
			}
		})
	}
}

func TestHandler_PageByID(t *testing.T) {
	h, repo := setupTestHandler()

	// Предзаполняем репозиторий для теста существующего ID
	targetURL := "https://yandex.ru"
	existingID := "shortKey123"
	_ = repo.Save(existingID, targetURL)

	tests := []struct {
		name       string
		urlID      string
		wantStatus int
		wantLoc    string
	}{
		{
			name:       "Редирект по существующему ID",
			urlID:      existingID,
			wantStatus: http.StatusTemporaryRedirect,
			wantLoc:    targetURL,
		},
		{
			name:       "ID не найден в базе",
			urlID:      "nonexistentID",
			wantStatus: http.StatusBadRequest,
			wantLoc:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get("/{id}", h.PageByID)

			request := httptest.NewRequest(http.MethodGet, "/"+tt.urlID, nil)
			recorder := httptest.NewRecorder()

			r.ServeHTTP(recorder, request)

			res := recorder.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantStatus, res.StatusCode)
			if tt.wantLoc != "" {
				assert.Equal(t, tt.wantLoc, res.Header.Get("Location"))
			}
		})
	}
}
