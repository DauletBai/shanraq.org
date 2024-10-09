package services

import (
	"errors"
	"shanraq.org/internal/models"
	"shanraq.org/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo repositories.UserRepository
	TokenService JWTService
}

func (s *AuthService) Register(user *models.User) error {
	// User data validation
	if err := user.Validate(); err != nil {
		return err
	}

	// Password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Save user database
	return &s.UserRepo.Create(user)
}

func (s *AuthService) Login(email, password string) (string error) {
	// Getting user by email
	user, err := s.UserRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("Incorred email or password")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("Incorred email or password")
	}

	// JWT token generation
	token, err := s.TokenService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}