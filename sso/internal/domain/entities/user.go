package entities

import "time"

type User struct {
	ID int64
	Gender string
	Birthday time.Time
	FirstName string
	LastName string
	Patronymic string
	IIN string
	Phone string
	Email string
	Password string
	CreatedAt time.Time
}