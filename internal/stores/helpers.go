// Package stores - пакет хранилищ
package stores

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/vbatushev/gopractice/internal/models"
)

// findContact возвращает индекс контакта с заданным ID в срезе контактов.
// Возвращает -1, если контакт не найден.
func findContact(contacts []models.Contact, id int64) int {
	return slices.IndexFunc(contacts, func(contact models.Contact) bool {
		return contact.ID == id
	})
}

// checkContentType проверяет, что Content-Type запроса соответствует ожидаемому медиа-типу.
// Возвращает код ответа и ошибку, если проверка не пройдена.
func checkContentType(r *http.Request, mediaType string) (int, error) {
	contentType := r.Header.Get("Content-Type")

	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if mediatype != mediaType {
		return http.StatusUnsupportedMediaType, fmt.Errorf(`Ожидается Content-Type=%s`, mediaType)
	}

	return http.StatusOK, nil
}

// validateErrors валидирует контакт и возвращает JSON-строку с ошибками валидации.
// Если ошибок нет, возвращает пустую строку и код 200 OK.
func validateErrors(contact models.Contact, vldr *validator.Validate) (string, int) {
	err := vldr.Struct(contact)
	if err != nil {
		var validateErrors models.Errors
		if v, ok := errors.AsType[validator.ValidationErrors](err); ok {
			for _, fr := range v {
				validateErrors.Errors = append(validateErrors.Errors, models.ValidationError{
					Field:  fr.Field(),
					Reason: fr.Error(),
				})
			}
		}
		resultParseValidate, errValidate := json.Marshal(validateErrors)
		if errValidate != nil {
			return fmt.Sprintf(`{ "status": "failed", "error": "%s" }`, errValidate.Error()), http.StatusBadRequest
		}
		return string(resultParseValidate), http.StatusBadRequest
	}
	return "", http.StatusOK
}

func responceJSONError(w http.ResponseWriter, status int, message string) error {
	respJSON := models.ResponceJSON{
		Message: message,
	}
	if status == http.StatusOK {
		respJSON.Status = "success"
	} else {
		respJSON.Status = "failed"
	}

	data, err := json.Marshal(respJSON)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(data)
	return err
}
