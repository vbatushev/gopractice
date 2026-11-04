// Package main - основной пакет
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
	mux := http.NewServeMux()

	// ENDPOINTS
	mux.HandleFunc(server.HealthPath, practice.HealthHandler)
	mux.HandleFunc(server.ContactsPath, practice.ContactsHandler)
	mux.HandleFunc(server.ContactPath, practice.ContactHandler)

	err := http.ListenAndServe(HTTPAddress, mux)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
