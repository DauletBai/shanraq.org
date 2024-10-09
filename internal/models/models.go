package models

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Role      string    `json:"role" db:"role"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdateAt  time.Time `json:"update_at" db:"update_at"`
}

func (u *User) Validate() error {
	if string.TrimSpace(u.Name) == "" {
		return errors.New("Name is required")
	}
	if string.TrimSpace(u.Email) == "" {
		return errors.New("Email is required")
	}
	if string.TrimSpace(u.Password) == "" {
		return errors.New("Password is required")
	}
	if u.Role != "admin" && u.Role != "editor" && u.Role != "author" && u.Role != "guest" {
		return errors.New("Invalid roles")
	}
	return nil
}

type Article struct {
	ID         int       `json:"id" db:"id"`
	Title      string    `json:"title" db:"title"`
	Content    string    `json:"content" db:"content"`
	AuthorID   int       `json:"aurhor_id" db:"author_id"`
	CategoryID int       `json:"category_id" db:"category_id"`
	CreatedAt  time.Time `json:"create_at" db:"create_at"`
	UpdateAt   time.Time `json:"update_at" db:"update_at"`
}

func (a *Article) Validate() error {
	if strings.TrimSpace(a.Title) == "" {
		return errors.New("Title is required")
	}
	if strings.TrimSpace(a.Content) == "" {
		return errors.New("Content is required")
	}
	if a.AuthorID <= 0 {
		return errors.New("Invalid author id")
	}
	if a.CategoryID <= 0 {
		return errors.New("Invalid category id")
	}
	return nil
}

// Similar validation for category
