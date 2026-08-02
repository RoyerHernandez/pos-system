package web

import (
	"context"

	pkgauth "github.com/RoyerHernandez/pos-system/api/pkg/infrastructure/auth"
)

// Claims represents the authenticated user's claims extracted from JWT.
type Claims struct {
	UserID int
	Role   string
}

// UserFromContext extracts auth claims from the request context.
func UserFromContext(ctx context.Context) *Claims {
	c := pkgauth.UserFromContext(ctx)
	if c == nil {
		return nil
	}
	return &Claims{
		UserID: c.UserID,
		Role:   c.Role,
	}
}
