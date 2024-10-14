package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Middleware для логирования запросов
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()           // Запоминаем время начала обработки запроса
		next(w, r)                    // Позволяем следующему обработчику выполнить свою работу
		duration := time.Since(start) // Вычисляем продолжительность
		log.Printf(
			"Method: %s, URL: %s, Duration: %s\n",
			r.Method,
			r.URL.Path,
			duration,
		)
	}
}

// Обработчик для приветственного сообщения
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hi!"))
}

// Обработчик для получения данных
func dataHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var data map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		fmt.Println("Data that we got:", data)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Got a data!"))

	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Get request /data"))

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	// Регистрация маршрутов с использованием middleware
	http.HandleFunc("/hello", loggingMiddleware(helloHandler))
	http.HandleFunc("/data", loggingMiddleware(dataHandler))

	// Запуск сервера
	fmt.Println("Running server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error occured while running server:", err)
	}
}

// Invoke-RestMethod -Uri http://localhost:8080/data -Method POST -Headers @{ "Content-Type" = "application/json" } -Body '{"name": "John", "age": 30}'
