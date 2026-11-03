// Package models - модуль описывающий структуры
package models

// Contact - модель контакта
type Contact struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	MiddleName  string `json:"middle_name"`
	LastName    string `json:"last_name"`
	Birthday    string `json:"birthday"`
	Description string `json:"description"`
}
