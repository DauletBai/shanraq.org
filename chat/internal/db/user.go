package db

import (
	"database/sql"

	"shanraq.org/chat/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := 
		INSERT INTO users (gender, birthday, first_name, last_name, patronymic, iin, phone, email, password, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW()) 
		RETURNING id, created_at 
	err := r.DB.QueryRow(query, 
		user.Gender,
		user.Birthday,
		user.FirstName,
		user.LastName,
		user.Patronymic,
		user.IIN,
		user.Phone,
		user.Email,
		user.Password,
	) .Scan(&user.ID, &user.CreatedAt)
	return err
}

// Add func: GetUserByEmail, GetUserByID ...