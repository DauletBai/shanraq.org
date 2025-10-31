package jobs

import (
	"net/http"

	"github.com/google/uuid"
	"shanraq.org/pkg/modules/auth"
)

// TenantResolver extracts the user ID from request context to scope job operations.
type TenantResolver func(r *http.Request) (uuid.UUID, bool)

// AuthTenantResolver returns a resolver that reads auth claims from the request context.
func AuthTenantResolver() TenantResolver {
	return func(r *http.Request) (uuid.UUID, bool) {
		claims, ok := auth.ClaimsFromContext(r.Context())
		if !ok {
			return uuid.Nil, false
		}
		id, err := uuid.Parse(claims.UserID)
		if err != nil {
			return uuid.Nil, false
		}
		return id, true
	}
}
