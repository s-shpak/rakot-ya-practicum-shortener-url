package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/rakot9/ya-practicum-shortener-url/internal/config"
	"net/http"
	"testing"
	"time"
)

func Test_run(t *testing.T) {

	r := chi.NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	// 2. Конфигурируем тестовые флаги (используем свободный порт localhost)
	testAddr := "localhost:8081"
	flags := config.Flags{
		FlagRunAddr:           testAddr,
		FlagRunShorternerAddr: "localhost:8081",
		FlagLog:               true,
	}

	// 3. Запускаем сервер в отдельной горутине, чтобы он не заблокировал выполнение теста
	errChan := make(chan error, 1)
	go func() {
		errChan <- run(r, flags)
	}()

	// Даем серверу небольшую задержку на запуск
	time.Sleep(50 * time.Millisecond)

	// 4. Выполняем реальный HTTP-запрос к запущенному тестовому серверу
	resp, err := http.Get("http://" + testAddr + "/ping")
	if err != nil {
		t.Fatalf("Не удалось отправить запрос к тестовому серверу: %v", err)
	}
	defer resp.Body.Close()

	// 5. Проверяем корректность ответа сервера
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200 OK, получен: %d", resp.StatusCode)
	}
}
