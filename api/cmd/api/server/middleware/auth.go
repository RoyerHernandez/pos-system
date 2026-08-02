package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	pkgauth "github.com/RoyerHernandez/pos-system/api/pkg/infrastructure/auth"
)

// AuthMiddleware validates JWT tokens from the Authorization header.
func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				respondUnauthorized(w, "missing authorization header")
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondUnauthorized(w, "invalid authorization format")
				return
			}

			claims, err := pkgauth.ValidateToken(parts[1], secret)
			if err != nil {
				respondUnauthorized(w, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), pkgauth.UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserFromContext extracts auth claims from the request context.
func UserFromContext(ctx context.Context) *pkgauth.Claims {
	return pkgauth.UserFromContext(ctx)
}

func respondUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"code":    "UNAUTHORIZED",
			"message": message,
		},
	})
}
