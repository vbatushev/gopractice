package endpoints

import (
	"log"
	"net/http"
)

const HealthEndpoint = "/api/v1/health"

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("HEALTHY"))
	if err != nil {
		log.Print(err)
	}
}
