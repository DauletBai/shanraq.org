package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTokenServiceGenerateAndParse(t *testing.T) {
	svc := NewTokenService("secret", time.Minute)
	user := User{ID: uuid.New(), Email: "alice@example.com"}

	token, err := svc.Generate(user)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	claims, err := svc.Parse(token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}

	if claims.UserID != user.ID.String() {
		t.Fatalf("expected user id %s, got %s", user.ID, claims.UserID)
	}
	if claims.Email != user.Email {
		t.Fatalf("expected email %s, got %s", user.Email, claims.Email)
	}
}

func TestTokenServiceDetectsWrongSecret(t *testing.T) {
	svc := NewTokenService("secret-a", time.Minute)
	user := User{ID: uuid.New(), Email: "bob@example.com"}

	token, err := svc.Generate(user)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	other := NewTokenService("secret-b", time.Minute)
	if _, err := other.Parse(token); err == nil {
		t.Fatalf("expected parse error with wrong secret")
	}
}
