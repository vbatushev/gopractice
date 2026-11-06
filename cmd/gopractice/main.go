// Package main - основной пакет
package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

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

	srv := http.Server{
		Handler:           mux,
		ReadTimeout:       45 * time.Second,
		WriteTimeout:      45 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ln, err := net.Listen("tcp", HTTPAddress)
	srv.Serve(ln)

	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
