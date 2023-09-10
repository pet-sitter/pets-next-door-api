package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	firebaseinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/firebase"
	s3infra "github.com/pet-sitter/pets-next-door-api/internal/infra/s3"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
	httpSwagger "github.com/swaggo/http-swagger/v2"
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
	r.Use(buildFirebaseAuthMiddleware(authClient))
}

func addRoutes(r *chi.Mux) {
	db, err := database.Open(configs.DatabaseURL)
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}

	authHandler := newAuthHandler()

	mediaService := media.NewMediaService(
		postgres.NewMediaPostgresStore(db),
		s3infra.NewS3Client(
			configs.B2KeyID,
			configs.B2Key,
			configs.B2Endpoint,
			configs.B2Region,
			configs.B2BucketName,
		),
	)
	mediaHandler := newMediaHandler(mediaService)

	userService := user.NewUserService(
		postgres.NewUserPostgresStore(db),
		postgres.NewPetPostgresStore(db),
		mediaService,
	)
	userHandler := newUserHandler(userService)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Get("/login/kakao", authHandler.kakaoLogin)
			r.Get("/callback/kakao", authHandler.kakaoCallback)
		})
		r.Route("/media", func(r chi.Router) {
			r.Get("/{id}", mediaHandler.findMediaByID)
			r.Post("/images", mediaHandler.uploadImage)
		})
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.RegisterUser)
			r.Post("/status", userHandler.FindUserStatusByEmail)
			r.Get("/me", userHandler.FindMyProfile)
			r.Put("/me", userHandler.UpdateMyProfile)
			r.Get("/me/pets", userHandler.FindMyPets)
			r.Put("/me/pets", userHandler.AddMyPets)
		})
	})
}
