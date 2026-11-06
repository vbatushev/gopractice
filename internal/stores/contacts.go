// Package stores - пакет хранилищ
package stores

import (
	"encoding/json/v2"
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

	contacts []models.Contact
	lastID   int64
	validate *validator.Validate
}

// NewContactStore - новое хранилище для контактов
func NewContactStore() *ContactStore {
	store := &ContactStore{}
	store.validate = validator.New()
	store.lastID = 1
	return store
}

// Add - создание контакта
func (store *ContactStore) Add(w http.ResponseWriter, r *http.Request) {
	store.Lock()
	defer store.Unlock()

	w.Header().Set("Content-Type", "application/json")

	status, err := checkContentType(r, "application/json")
	if err != nil {
		responceJSONError(w, status, err.Error())
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var contact models.Contact
	err = json.UnmarshalRead(r.Body, &contact, json.DefaultOptionsV2())

	if err != nil {
		if _, ok := err.(*http.MaxBytesError); ok {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	contact.ID = store.lastID

	errorsString, status := validateErrors(contact, store.validate)
	if status != http.StatusOK {
		http.Error(w, errorsString, status)
		return
	}

	store.contacts = append(store.contacts, contact)
	store.lastID++

	result, err := json.Marshal(contact)
	if err != nil {
		responceJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(result))
}

// GetAll - получение пользователей
func (store *ContactStore) GetAll(w http.ResponseWriter, _ *http.Request) {
	store.Lock()
	defer store.Unlock()

	contacts, err := json.Marshal(store.contacts)
	if err != nil {
		responceJSONError(w, http.StatusInternalServerError, err.Error())
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

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		responceJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	indx := findContact(store.contacts, int64(id))
	if indx < 0 {
		responceJSONError(w, http.StatusNotFound, "не найден пользователь")
		return
	}

	result, err := json.Marshal(store.contacts[indx])
	if err != nil {
		responceJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Write(result)
}

// Remove - удаление контакта
func (store *ContactStore) Remove(w http.ResponseWriter, r *http.Request) {
	store.Lock()
	defer store.Unlock()

	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		responceJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	indx := findContact(store.contacts, int64(id))
	if indx < 0 {
		responceJSONError(w, http.StatusNotFound, "не найден пользователь")
		return
	}

	store.contacts = slices.Delete(store.contacts, indx, indx+1)

	w.WriteHeader(http.StatusNoContent)
}

// GetContacts - getter для конактов
func (store *ContactStore) GetContacts() []models.Contact {
	return store.contacts
}
