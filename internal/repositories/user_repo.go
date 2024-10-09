package repositories

import (
	"database/sql"
	"shanraq.org/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	// Other methods
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (name, email, password, role, create_at, update_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id`
	return r.db.QueryRow(query, user.Name, user.Email, user.Password, user.Role) .Scan(&user.ID)
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, email, password, role FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}