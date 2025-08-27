package models

import "github.com/google/uuid"

type Country struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Population int64     `json:"population"`
	FlagCode   string    `json:"flag_code"`
}

// ... в будущем здесь будут и другие модели: Category, Competition, Award ...