package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Приветственное сообщение
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hi!"))
}

// Обработчик POST запросов с данными в формате JSON
func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Вывод данных в консоль
	fmt.Println("Got a data:", data)

	// Ответ клиенту
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Got a data!"))
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/data", dataHandler)

	fmt.Println("Running a server in port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error occured while running server:", err)
	}
}

// Invoke-RestMethod -Uri http://localhost:8080/data -Method POST -Headers @{ "Content-Type" = "application/json" } -Body '{"name": "John", "age": 30}'
