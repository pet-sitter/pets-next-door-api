package main

import (
	"context"
	"encoding/json"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pet-sitter/pets-next-door-api/cmd/server/handler"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/auth"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	firebaseinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/firebase"
	kakaoinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/kakao"
	s3infra "github.com/pet-sitter/pets-next-door-api/internal/infra/s3"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
	pndMiddleware "github.com/pet-sitter/pets-next-door-api/lib/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(app *firebaseinfra.FirebaseApp) *chi.Mux {
	r := chi.NewRouter()

	db, err := database.Open(configs.DatabaseURL)
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}

	authClient, err := app.Auth(context.Background())

	// Initialize services
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

	userService := user.NewUserService(
		postgres.NewUserPostgresStore(db),
		postgres.NewPetPostgresStore(db),
		mediaService,
	)

	authService := auth.NewFirebaseBearerAuthService(authClient, userService)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	breedService := pet.NewBreedService(
		postgres.NewBreedPostgresStore(db),
	)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, kakaoinfra.NewKakaoClient())
	userHandler := handler.NewUserHandler(userService, authService)
	mediaHandler := handler.NewMediaHandler(mediaService)
	breedHandler := handler.NewBreedHandler(breedService)

	// Register middlewares
	r.Use(middleware.Logger)
	r.Use(pndMiddleware.BuildAuthMiddleware(authService, auth.FirebaseAuthClientKey))

	// Register routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Get("/login/kakao", authHandler.KakaoLogin)
			r.Get("/callback/kakao", authHandler.KakaoCallback)
		})
		r.Route("/media", func(r chi.Router) {
			r.Get("/{id}", mediaHandler.FindMediaByID)
			r.Post("/images", mediaHandler.UploadImage)
		})
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.RegisterUser)
			r.Post("/check/nickname", userHandler.CheckUserNickname)
			r.Post("/status", userHandler.FindUserStatusByEmail)
			r.Get("/me", userHandler.FindMyProfile)
			r.Put("/me", userHandler.UpdateMyProfile)
			r.Get("/me/pets", userHandler.FindMyPets)
			r.Put("/me/pets", userHandler.AddMyPets)
		})
		r.Route("/breeds", func(r chi.Router) {
			r.Get("/", breedHandler.FindBreeds)
		})
	})

	return r
}
