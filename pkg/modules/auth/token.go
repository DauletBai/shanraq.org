package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenService issues and validates JWT tokens.
type TokenService struct {
	secret []byte
	ttl    time.Duration
}

// Claims wraps jwt.RegisteredClaims for convenience.
type Claims struct {
	UserID      string   `json:"uid"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	PrimaryRole string   `json:"role,omitempty"`
	jwt.RegisteredClaims
}

func NewTokenService(secret string, ttl time.Duration) *TokenService {
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	return &TokenService{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (s *TokenService) Generate(user User) (string, error) {
	primary, roles := normalizeClaimRoles(user)
	claims := Claims{
		UserID:      user.ID.String(),
		Email:       user.Email,
		Roles:       roles,
		PrimaryRole: primary,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "shanraq",
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

func (s *TokenService) TTL() time.Duration {
	return s.ttl
}

func (s *TokenService) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		claims.normalize()
		return claims, nil
	}
	return nil, errors.New("invalid token claims")
}

func (c *Claims) normalize() {
	combined := append([]string{c.PrimaryRole}, c.Roles...)
	combined = sanitizeRoles(combined)
	if len(combined) == 0 {
		combined = []string{defaultRoleName}
	}
	c.PrimaryRole = combined[0]
	c.Roles = combined
}

// HasAnyRole returns true when the claim includes any of the provided roles.
func (c *Claims) HasAnyRole(roles ...string) bool {
	if len(roles) == 0 {
		return true
	}
	c.normalize()
	allowed := sanitizeRoles(roles)
	if len(allowed) == 0 {
		return true
	}
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, role := range allowed {
		allowedSet[role] = struct{}{}
	}
	for _, role := range c.Roles {
		if _, ok := allowedSet[role]; ok {
			return true
		}
	}
	return false
}

func normalizeClaimRoles(user User) (string, []string) {
	combined := append([]string{user.Role}, user.Roles...)
	combined = append(combined, defaultRoleName)
	normalized := sanitizeRoles(combined)
	if len(normalized) == 0 {
		normalized = []string{defaultRoleName}
	}
	return normalized[0], normalized
}

// ClaimsForUser builds a Claims struct for the provided user without signing a token.
func ClaimsForUser(user User) *Claims {
	primary, roles := normalizeClaimRoles(user)
	return &Claims{
		UserID:      user.ID.String(),
		Email:       user.Email,
		Roles:       roles,
		PrimaryRole: primary,
	}
}
