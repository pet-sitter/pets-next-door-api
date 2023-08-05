package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/database"
	firebaseinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/firebase"
	s3infra "github.com/pet-sitter/pets-next-door-api/internal/infra/s3"
	"github.com/pet-sitter/pets-next-door-api/internal/media"
	"github.com/pet-sitter/pets-next-door-api/internal/user"
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

	db.Migrate(configs.MigrationPath)

	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}

	userService := user.NewUserService(db)
	userHandler := newUserHandler(userService)
	authHandler := newAuthHandler()

	s3Client := s3infra.NewS3Client(
		configs.B2KeyID,
		configs.B2Key,
		configs.B2Endpoint,
		configs.B2Region,
		configs.B2BucketName,
	)
	mediaService := media.NewMediaService(db, s3Client)
	mediaHandler := newMediaHandler(mediaService)

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
		})
	})
}
