// Package stores - пакет хранилищ
package stores

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"slices"
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
	if indx < 0 {
		contact.ID = store.getLastID()
		store.Contacts = append(store.Contacts, contact)
	} else {
		store.Contacts[indx] = contact
	}

	w.Write([]byte(`{ "status": "success", error: null }`))
}

// GetAll - получение пользователей
func (store *ContactStore) GetAll(w http.ResponseWriter, r *http.Request) ([]models.Contact, error) {
	store.Lock()
	defer store.Unlock()

	r.Body.Close()

	w.Write([]byte("test"))

	return store.Contacts, nil
}

// Get - получение контакта по ID
func (store *ContactStore) Get(id int) (models.Contact, error) {
	store.Lock()
	defer store.Unlock()

	indx := findUser(store.Contacts, id)
	if indx < 0 {
		return models.Contact{}, errors.New("Не нашлось ничего")
	}

	return store.Contacts[indx], nil
}

// Remove - удаление контакта
func (store *ContactStore) Remove(w http.ResponseWriter, r *http.Request) error {
	store.Lock()
	defer store.Unlock()

	r.Body.Close()

	w.Write([]byte("test"))

	return nil
}

func (store *ContactStore) getLastID() int {
	if len(store.Contacts) == 0 {
		return 0
	}

	maxID := 0
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
