package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	registerMiddlewares(r)
	addRoutes(r)

	return r
}

func registerMiddlewares(r *chi.Mux) {
	r.Use(middleware.Logger)
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
