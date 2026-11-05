// Package stores - пакет хранилищ
package stores

import (
	"encoding/json/v2"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/vbatushev/gopractice/internal/models"
)

// ContactStore - структура хранилища контактов
type ContactStore struct {
	sync.Mutex

	Contacts []models.Contact
	lastID   int64
	validate *validator.Validate
}

// NewContactStore - новое хранилище для контактов
func NewContactStore() *ContactStore {
	store := &ContactStore{}
	store.validate = validator.New()
	return store
}

// Add - создание контакта
func (store *ContactStore) Add(w http.ResponseWriter, r *http.Request) {
	store.Lock()
	defer store.Unlock()

	w.Header().Set("Content-Type", "application/json")

	status, err := checkContentType(r, "application/json")
	if err != nil {
		http.Error(w, fmt.Sprintf(`{ "status": "failed", "error": "%s" }`, err.Error()), status)
	}

	var contact models.Contact

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{ "status": "failed", "error": "Не удалость прочитать тело запроса" }`, http.StatusBadRequest)
	}

	err = json.Unmarshal(data, &contact, json.DefaultOptionsV2())
	if err != nil {
		http.Error(w, fmt.Sprintf(`{ "status": "failed", "error": "%s" }`, err.Error()), http.StatusBadRequest)
		return
	}

	contact.ID = store.lastID

	errorsString, status := validateErrors(contact, store.validate)
	if status != http.StatusOK {
		http.Error(w, errorsString, status)
	}

	store.Contacts = append(store.Contacts, contact)
	store.lastID++

	result, err := json.Marshal(contact)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{ "status": "failed", "error": "%s" }`, err.Error()), http.StatusBadRequest)
		return
	}

	w.Write([]byte(result))
}

// GetAll - получение пользователей
func (store *ContactStore) GetAll(w http.ResponseWriter, _ *http.Request) {
	store.Lock()
	defer store.Unlock()

	contacts, err := json.Marshal(store.Contacts)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{ "status": "failed", "error": "%s" }`, err.Error()), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(contacts)
}

// Get - получение контакта по ID
func (store *ContactStore) Get(w http.ResponseWriter, r *http.Request) {
	store.Lock()
	defer store.Unlock()

	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, `{ "status": "failed", "error": "неправильный id" }`, http.StatusBadRequest)
		return
	}

	indx := findContact(store.Contacts, int64(id))
	if indx < 0 {
		http.Error(w, `{ "status": "failed", "error": "не найден пользователь" }`, http.StatusNotFound)
		return
	}

	result, err := json.Marshal(store.Contacts[indx])
	if err != nil {
		http.Error(w, fmt.Sprintf(`{ "status": "failed", "error": "%s" }`, err.Error()), http.StatusNotFound)
		return
	}

	w.Write(result)
}

// Remove - удаление контакта
func (store *ContactStore) Remove(w http.ResponseWriter, r *http.Request) {
	store.Lock()
	defer store.Unlock()

	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, `{ "status": "failed", "error": "неправильный id" }`, http.StatusBadRequest)
		return
	}

	indx := findContact(store.Contacts, int64(id))
	if indx < 0 {
		http.Error(w, `{ "status": "failed", "error": "не найден пользователь" }`, http.StatusNotFound)
		return
	}

	store.Contacts = slices.Delete(store.Contacts, indx, indx+1)

	w.WriteHeader(http.StatusNoContent)
}
