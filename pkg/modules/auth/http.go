package auth

import (
	"errors"
	"strings"
)

type signupRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
	// Consent to the Terms and Privacy Policy is mandatory (KZ online-platform
	// law); the API must not be a way around the browser consent checkbox.
	Consent bool `json:"consent" form:"consent"`
	// A real first and last name, on the same terms as the browser form. Every
	// article, listing and comment on the site is attributed to a person, so an
	// account with no name is an account that cannot be attributed — and the API
	// used to create exactly that, quietly making the real-name rule optional
	// for anyone who preferred JSON to the form.
	FirstName  string `json:"first_name" form:"first_name"`
	LastName   string `json:"last_name" form:"last_name"`
	MiddleName string `json:"middle_name" form:"middle_name"`
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
