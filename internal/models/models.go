// Package models - модуль описывающий структуры
package models

// Contact - модель контакта
type Contact struct {
	ID          int64  `json:"id"`
	FirstName   string `json:"first_name" validate:"required,min=1,max=255"`
	MiddleName  string `json:"middle_name" validate:"max=255"`
	LastName    string `json:"last_name" validate:"required,min=1,max=255"`
	Birthday    string `json:"birthday" validate:"datetime=2006-01-02"`
	Description string `json:"description" validate:"max=255"`
}

// ValidationError - ошибка валидации
type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

// Errors - ошибки валидации
type Errors struct {
	Errors []ValidationError
}
