// Package stores - пакет хранилищ
package stores

import (
	"encoding/json/v2"
	"fmt"
	"io"
	"mime"
	"net/http"
	"slices"
	"strconv"
	"sync"

	"github.com/vbatushev/gopractice/internal/models"
)

// ContactStore - структура хранилища контактов
type ContactStore struct {
	sync.Mutex

	Contacts []models.Contact
}

// NewContactStore - новое хранилище для контактов
func NewContactStore() *ContactStore {
	store := &ContactStore{}
	return store
}

// Add - создание контакта
func (store *ContactStore) Add(w http.ResponseWriter, r *http.Request) {
	store.Lock()
	defer store.Unlock()

	contentType := r.Header.Get("Content-Type")
	w.Header().Set("Content-Type", "application/json")

	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{ "status": "failed", "error": "%s" }`, err.Error()), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		http.Error(w, `{ "status": "failed", "error": "Ожидается Content-Type=application/json" }`, http.StatusUnsupportedMediaType)
		return
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

	indx := findUser(store.Contacts, contact.ID)
	if indx < 1 {
		contact.ID = store.getLastID()
		store.Contacts = append(store.Contacts, contact)
	} else {
		store.Contacts[indx] = contact
	}

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

	indx := findUser(store.Contacts, id)
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

	indx := findUser(store.Contacts, id)
	if indx < 0 {
		http.Error(w, `{ "status": "failed", "error": "не найден пользователь" }`, http.StatusNotFound)
		return
	}

	store.Contacts = slices.Delete(store.Contacts, indx, indx+1)

	w.WriteHeader(http.StatusNoContent)
}

func (store *ContactStore) getLastID() int {
	if len(store.Contacts) == 0 {
		return 1
	}

	maxID := 1
	for _, cnt := range store.Contacts {
		if cnt.ID > maxID {
			maxID = cnt.ID
		}
	}
	return maxID + 1
}

func findUser(contacts []models.Contact, id int) int {
	return slices.IndexFunc(contacts, func(contact models.Contact) bool {
		return contact.ID == id
	})
}
