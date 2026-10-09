package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter(t *testing.T) {
	r := Router("localhost:8080", "http://localhost:8080", false)

	// Структура табличного теста
	tests := []struct {
		name           string
		method         string
		url            string
		body           string
		expectedStatus int
	}{
		{
			name:           "POST - create short URL",
			method:         http.MethodPost,
			url:            "/",
			body:           "https://yandex.ru",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "GET - root URL with wrong method",
			method:         http.MethodGet,
			url:            "/",
			body:           "",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Создаем виртуальный HTTP-запрос
			req, err := http.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", "text/plain")

			// 2. Создаем инструмент для записи ответа
			rec := httptest.NewRecorder()

			// 3. Отправляем запрос напрямую в тестируемый роутер
			r.ServeHTTP(rec, req)

			// 4. Проверяем соответствие статус-кода ответа
			if rec.Code != tc.expectedStatus {
				t.Errorf("case '%s' failed: expected status %d, got %d", tc.name, tc.expectedStatus, rec.Code)
			}
		})
	}
}
