package ports

import "shanraq.org/sso/internal/domain/entities"

type UserRepository interface {
	Create(user *entities.User) error
	GetByEmail(email string) (*entities.User, error)
	GetByID(id int64) (*entities.User, error)
}

