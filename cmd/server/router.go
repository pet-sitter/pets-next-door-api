package main

import (
	"context"
	"log"
	"net/http"

	"github.com/pet-sitter/pets-next-door-api/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/pet-sitter/pets-next-door-api/cmd/server/handler"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/auth"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	firebaseinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/firebase"
	kakaoinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/kakao"
	s3infra "github.com/pet-sitter/pets-next-door-api/internal/infra/s3"
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
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Initialize services
	s3Client := s3infra.NewS3Client(
		configs.B2KeyID,
		configs.B2Key,
		configs.B2Endpoint,
		configs.B2Region,
		configs.B2BucketName,
	)

	mediaService := service.NewMediaService(db, s3Client)
	userService := service.NewUserService(db, s3Client)
	authService := service.NewFirebaseBearerAuthService(db, authClient, s3Client)
	userDeleter := service.NewUserDeleter(db, authService)
	breedService := service.NewBreedService(db)
	sosPostService := service.NewSosPostService(db)
	conditionService := service.NewConditionService(db)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, kakaoinfra.NewKakaoDefaultClient())
	userHandler := handler.NewUserHandler(*userService, authService)
	userDeleterHandler := handler.NewUserDeleterHandler(*userDeleter, authService)
	mediaHandler := handler.NewMediaHandler(*mediaService)
	breedHandler := handler.NewBreedHandler(*breedService)
	sosPostHandler := handler.NewSosPostHandler(*sosPostService, authService)
	conditionHandler := handler.NewConditionHandler(*conditionService)

	// Register middlewares
	r.Use(middleware.Logger)
	r.Use(pndMiddleware.BuildAuthMiddleware(authService, auth.FirebaseAuthClientKey))

	// Register routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]string{"status": "ok"})
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Get("/login/kakao", authHandler.KakaoLogin)
			r.Get("/callback/kakao", authHandler.KakaoCallback)
			r.Post("/custom-tokens/kakao", authHandler.GenerateFBCustomTokenFromKakao)
		})
		r.Route("/media", func(r chi.Router) {
			r.Get("/{id}", mediaHandler.FindMediaByID)
			r.Post("/images", mediaHandler.UploadImage)
		})
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.RegisterUser)
			r.Post("/check/nickname", userHandler.CheckUserNickname)
			r.Post("/status", userHandler.FindUserStatusByEmail)
			r.Get("/", userHandler.FindUsers)
			r.Get("/me", userHandler.FindMyProfile)
			r.Put("/me", userHandler.UpdateMyProfile)
			r.Get("/me/pets", userHandler.FindMyPets)
			r.Put("/me/pets", userHandler.AddMyPets)
			r.Delete("/me/hard", userDeleterHandler.HardDeleteMyAccount)
		})
		r.Route("/breeds", func(r chi.Router) {
			r.Get("/", breedHandler.FindBreeds)
		})
		r.Route("/posts", func(r chi.Router) {
			r.Post("/sos", sosPostHandler.WriteSosPost)
			r.Get("/sos/{id}", sosPostHandler.FindSosPostByID)
			r.Get("/sos", sosPostHandler.FindSosPosts)
			r.Put("/sos", sosPostHandler.UpdateSosPost)
			r.Get("/sos/conditions", conditionHandler.FindConditions)
		})
	})

	return r
}
