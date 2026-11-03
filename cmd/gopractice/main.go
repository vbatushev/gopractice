package main

import (
	"fmt"
	"net/http"

	"github.com/vbatushev/gopractice/internal/endpoints"
)

const (
	HttpAddress = "127.0.0.1:8080"
)

func main() {
	// ENDPOINTS
	http.HandleFunc(endpoints.HealthEndpoint, endpoints.HealthHandler)

	err := http.ListenAndServe(HttpAddress, nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
