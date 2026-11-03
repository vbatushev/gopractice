package main

import (
	"fmt"
	"net/http"

	"github.com/vbatushev/gopractice/internal/server"
)

const (
	// HTTPAddress - привязка сервера
	HTTPAddress = "127.0.0.1:8080"
)

func main() {
	practice := server.NewPracticeServer()

	// ENDPOINTS
	http.HandleFunc(server.HealthPath, practice.HealthHandler)
	http.HandleFunc(server.ContactsPath, practice.ContactsHandler)

	err := http.ListenAndServe(HTTPAddress, nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
