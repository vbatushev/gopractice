// Package server - сервер
package server

import (
	"log"
	"net/http"

	"github.com/vbatushev/gopractice/internal/stores"
)

const (
	// HealthPath - эндпоинт контактов
	HealthPath = "/api/v1/health"
	// ContactsPath - эндпоинт контактов
	ContactsPath = "/api/v1/contacts"
	// ContactPath - эндпоинт контакта
	ContactPath = "/api/v1/contacts/{id}"
)

// PracticeServer - сервер практической работы
type PracticeServer struct {
	ContactStore *stores.ContactStore
}

// NewPracticeServer - создание сервера
func NewPracticeServer() *PracticeServer {
	contacts := stores.NewContactStore()
	return &PracticeServer{ContactStore: contacts}
}

// HealthHandler - проверяет работоспособность сервиса
func (srv *PracticeServer) HealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("HEALTHY"))
	if err != nil {
		log.Print(err)
	}
}

// ContactsHandler - обработчик API для контактов
func (srv *PracticeServer) ContactsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		srv.ContactStore.Add(w, r)
	case http.MethodGet:
		srv.ContactStore.GetAll(w, r)
	}
}

// ContactHandler - обработчик API для контакта
func (srv *PracticeServer) ContactHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		srv.ContactStore.Get(w, r)
	case http.MethodDelete:
		srv.ContactStore.Remove(w, r)
	}
}
