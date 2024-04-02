package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pet-sitter/pets-next-door-api/cmd/server/handler"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/auth"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	kakaoinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/kakao"
	s3infra "github.com/pet-sitter/pets-next-door-api/internal/infra/s3"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
	pndMiddleware "github.com/pet-sitter/pets-next-door-api/lib/middleware"
	"github.com/rs/zerolog"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"net/http"
	"os"

	firebaseinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/firebase"
)

func NewRouter(app *firebaseinfra.FirebaseApp) *echo.Echo {
	e := echo.New()
	ctx := context.Background()

	db, err := database.Open(configs.DatabaseURL)
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}

	authClient, err := app.Auth(ctx)
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
	userService := service.NewUserService(db, mediaService)
	authService := service.NewFirebaseBearerAuthService(authClient, userService)
	breedService := service.NewBreedService(db)
	sosPostService := service.NewSosPostService(db)
	conditionService := service.NewConditionService(db)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, kakaoinfra.NewKakaoDefaultClient())
	userHandler := handler.NewUserHandler(*userService, authService)
	mediaHandler := handler.NewMediaHandler(*mediaService)
	breedHandler := handler.NewBreedHandler(*breedService)
	sosPostHandler := handler.NewSosPostHandler(*sosPostService, authService)
	conditionHandler := handler.NewConditionHandler(*conditionService)

	// Register middlewares
	logger := zerolog.New(os.Stdout)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))
	e.Use(pndMiddleware.BuildAuthMiddleware(authService, auth.FirebaseAuthClientKey))

	// Register routes
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	apiRouteGroup := e.Group("/api")

	authApiGroup := apiRouteGroup.Group("/auth")
	{
		authApiGroup.GET("/login/kakao", authHandler.KakaoLogin)
		authApiGroup.GET("/callback/kakao", authHandler.KakaoCallback)
		authApiGroup.POST("/custom-tokens/kakao", authHandler.GenerateFBCustomTokenFromKakao)
	}

	mediaApiGroup := apiRouteGroup.Group("/media")
	{
		mediaApiGroup.GET("/:id", mediaHandler.FindMediaByID)
		mediaApiGroup.POST("/images", mediaHandler.UploadImage)
	}

	userApiGroup := apiRouteGroup.Group("/users")
	{
		userApiGroup.POST("", userHandler.RegisterUser)
		userApiGroup.POST("/check/nickname", userHandler.CheckUserNickname)
		userApiGroup.POST("/status", userHandler.FindUserStatusByEmail)
		userApiGroup.GET("", userHandler.FindUsers)
		userApiGroup.GET("/me", userHandler.FindMyProfile)
		userApiGroup.PUT("/me", userHandler.UpdateMyProfile)
		userApiGroup.DELETE("/me", userHandler.DeleteMyAccount)
		userApiGroup.GET("/me/pets", userHandler.FindMyPets)
		userApiGroup.PUT("/me/pets", userHandler.AddMyPets)
	}

	breedApiGroup := apiRouteGroup.Group("/breeds")
	{
		breedApiGroup.GET("", breedHandler.FindBreeds)
	}

	postApiGroup := apiRouteGroup.Group("/posts")
	{
		postApiGroup.POST("/sos", sosPostHandler.WriteSosPost)
		postApiGroup.GET("/sos/{id}", sosPostHandler.FindSosPostByID)
		postApiGroup.GET("/sos", sosPostHandler.FindSosPosts)
		postApiGroup.PUT("/sos", sosPostHandler.UpdateSosPost)
		postApiGroup.GET("/sos/conditions", conditionHandler.FindConditions)
	}

	return e
}
