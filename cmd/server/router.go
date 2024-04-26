package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pet-sitter/pets-next-door-api/cmd/server/handler"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/auth"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	kakaoinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/kakao"
	s3infra "github.com/pet-sitter/pets-next-door-api/internal/infra/s3"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
	pndmiddleware "github.com/pet-sitter/pets-next-door-api/lib/middleware"
	"github.com/rs/zerolog"
	echoswagger "github.com/swaggo/echo-swagger"

	firebaseinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/firebase"
)

func NewRouter(app *firebaseinfra.FirebaseApp) (*echo.Echo, error) {
	e := echo.New()
	ctx := context.Background()

	db, err := database.Open(configs.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %w", err)
	}

	// Initialize services
	s3Client, err := s3infra.NewS3Client(
		configs.B2KeyID,
		configs.B2Key,
		configs.B2Endpoint,
		configs.B2Region,
		configs.B2BucketName,
	)
	if err != nil {
		return nil, fmt.Errorf("error initializing s3 client: %w", err)
	}

	mediaService := service.NewMediaService(db, s3Client)
	userService := service.NewUserService(db, mediaService)
	authService := service.NewFirebaseBearerAuthService(authClient, userService)
	breedService := service.NewBreedService(db)
	sosPostService := service.NewSOSPostService(db)
	conditionService := service.NewSOSConditionService(db)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, kakaoinfra.NewKakaoDefaultClient())
	userHandler := handler.NewUserHandler(*userService, authService)
	mediaHandler := handler.NewMediaHandler(*mediaService)
	breedHandler := handler.NewBreedHandler(*breedService)
	sosPostHandler := handler.NewSOSPostHandler(*sosPostService, authService)
	conditionHandler := handler.NewConditionHandler(*conditionService)

	// Register middlewares
	logger := zerolog.New(os.Stdout)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))
	e.Use(pndmiddleware.BuildAuthMiddleware(authService, auth.FirebaseAuthClientKey))

	// Register routes
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	e.GET("/swagger/*", echoswagger.WrapHandler)

	apiRouteGroup := e.Group("/api")

	authAPIGroup := apiRouteGroup.Group("/auth")
	{
		authAPIGroup.GET("/login/kakao", authHandler.KakaoLogin)
		authAPIGroup.GET("/callback/kakao", authHandler.KakaoCallback)
		authAPIGroup.POST("/custom-tokens/kakao", authHandler.GenerateFBCustomTokenFromKakao)
	}

	mediaAPIGroup := apiRouteGroup.Group("/media")
	{
		mediaAPIGroup.GET("/:id", mediaHandler.FindMediaByID)
		mediaAPIGroup.POST("/images", mediaHandler.UploadImage)
	}

	userAPIGroup := apiRouteGroup.Group("/users")
	{
		userAPIGroup.POST("", userHandler.RegisterUser)
		userAPIGroup.POST("/check/nickname", userHandler.CheckUserNickname)
		userAPIGroup.POST("/status", userHandler.FindUserStatusByEmail)
		userAPIGroup.GET("", userHandler.FindUsers)
		userAPIGroup.GET("/me", userHandler.FindMyProfile)
		userAPIGroup.PUT("/me", userHandler.UpdateMyProfile)
		userAPIGroup.DELETE("/me", userHandler.DeleteMyAccount)
		userAPIGroup.GET("/me/pets", userHandler.FindMyPets)
		userAPIGroup.PUT("/me/pets", userHandler.AddMyPets)
		userAPIGroup.PUT("/me/pets/:petID", userHandler.UpdateMyPet)
		userAPIGroup.DELETE("/me/pets/:petID", userHandler.DeleteMyPet)
	}

	breedAPIGroup := apiRouteGroup.Group("/breeds")
	{
		breedAPIGroup.GET("", breedHandler.FindBreeds)
	}

	postAPIGroup := apiRouteGroup.Group("/posts")
	{
		postAPIGroup.POST("/sos", sosPostHandler.WriteSOSPost)
		postAPIGroup.GET("/sos/:id", sosPostHandler.FindSOSPostByID)
		postAPIGroup.GET("/sos", sosPostHandler.FindSOSPosts)
		postAPIGroup.PUT("/sos", sosPostHandler.UpdateSOSPost)
		postAPIGroup.GET("/sos/conditions", conditionHandler.FindConditions)
	}

	return e, nil
}
