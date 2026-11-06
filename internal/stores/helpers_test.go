package stores

import (
	"testing"

	"github.com/vbatushev/gopractice/internal/models"
)

func TestFindContact(t *testing.T) {
	contacts := []models.Contact{{
		ID:          1,
		FirstName:   "Александр",
		MiddleName:  "Сергеевич",
		LastName:    "Пушкин",
		Birthday:    "1799-06-06",
		Description: "Великий поэт",
	}, {
		ID:          3,
		FirstName:   "Кот",
		LastName:    "Матроскин",
		Birthday:    "1984-08-10",
		Description: "Кот из мультфильма",
	}, {
		ID:          2,
		FirstName:   "Пес",
		LastName:    "Шарик",
		Birthday:    "1984-11-1",
		Description: "Пес из Простоквашино",
	}}

	sharik := findContact(contacts, 2)
	if sharik != 2 {
		t.Errorf("Ожидалось, что Шарик будет третьим")
	}

	notfound := findContact(contacts, 4)
	if notfound != -1 {
		t.Errorf("Ожидалось, что при отсутствии контакта вернется -1")
	}
}
