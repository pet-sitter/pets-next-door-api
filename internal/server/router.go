package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	firebaseinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/firebase"
)

func NewRouter(app *firebaseinfra.FirebaseApp) *chi.Mux {
	r := chi.NewRouter()

	registerMiddlewares(r, app)
	addRoutes(r)

	return r
}

func registerMiddlewares(r *chi.Mux, app *firebaseinfra.FirebaseApp) {
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	r.Use(middleware.Logger)
	r.Use(buildFirebaseAppMiddleware(authClient))
}

func buildFirebaseAppMiddleware(app *auth.Client) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), firebaseAuthClientKey, app)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func addRoutes(r *chi.Mux) {
	authHandler := newAuthHandler()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Get("/login/kakao", authHandler.kakaoLogin)
			r.Get("/callback/kakao", authHandler.kakaoCallback)
		})
	})
}
