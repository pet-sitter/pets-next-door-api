package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/auth"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

func BuildAuthMiddleware(app service.AuthService, authKey auth.ContextKey) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), authKey, app)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
