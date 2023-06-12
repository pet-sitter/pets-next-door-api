package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
)

func buildFirebaseAuthMiddleware(app *auth.Client) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), firebaseAuthClientKey, app)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func verifyAuth(ctx context.Context, authHeader string) (*auth.Token, error) {
	authClient := ctx.Value(firebaseAuthClientKey).(*auth.Client)
	idToken, err := stripBearerToken(authHeader)
	if err != nil {
		return nil, err
	}

	authToken, err := authClient.VerifyIDToken(ctx, idToken)
	return authToken, err
}

func stripBearerToken(authHeader string) (string, error) {
	if len(authHeader) > 6 && strings.ToUpper(authHeader[0:7]) == "BEARER " {
		return authHeader[7:], nil
	}

	return authHeader, fmt.Errorf("invalid auth header")
}
