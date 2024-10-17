package services

import (
	"shanraq.org/sso/internal/application/dto"
	"shanraq.org/sso/internal/domain/entities"
	"shanraq.org/sso/internal/domain/ports"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo ports.UserRepository
	Auth     ports.Auth
}

func (s *UserService) Register(userDTO *dto.userDTO) error {
	// Data validation
	// ...

	// Password heshing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Creating a user entity
	user := &entities.User{
		Gender:     userDTO.Gender,
		Birthday:   userDTO.Birthday,
		FirstName:  userDTO.FirstName,
		LastName:   userDTO.LastName,
		Patronymic: userDTO.Patronymic,
		IIN:        userDTO.IIN,
		Phone:      userDTO.Phone,
		Email:      userDTO.Email,
		Password:   userDTO.Password,
	}

	// Save user
	return s.UserRepo.Create(user)
}

func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.UserRepo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	// Password check
	if err := bcrypt.CompareHashPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}

	// JWT token generation
	token, err := s.Auth.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
