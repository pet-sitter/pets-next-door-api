package middleware

import (
	"context"
	"net/http"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/auth"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

func BuildAuthMiddleware(app service.AuthService, authKey auth.ContextKey) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), authKey, app)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
