package auth

import (
	"errors"
	"strings"
)

type signupRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
}

type signinRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type passwordResetRequest struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}

type passwordResetConfirmRequest struct {
	Token    string `json:"token" form:"token" validate:"required"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
}

func tokenFromHeader(header string) (string, error) {
	if header == "" {
		return "", errors.New("missing authorization header")
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", errors.New("invalid authorization header")
	}
	return parts[1], nil
}
