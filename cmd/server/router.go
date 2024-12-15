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
	s3infra "github.com/pet-sitter/pets-next-door-api/internal/infra/bucket"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	kakaoinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/kakao"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
	"github.com/pet-sitter/pets-next-door-api/internal/wschat"
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
	chatService := service.NewChatService(db)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, kakaoinfra.NewKakaoDefaultClient())
	userHandler := handler.NewUserHandler(*userService, authService)
	mediaHandler := handler.NewMediaHandler(*mediaService)
	breedHandler := handler.NewBreedHandler(*breedService)
	sosPostHandler := handler.NewSOSPostHandler(*sosPostService, authService)
	conditionHandler := handler.NewConditionHandler(*conditionService)
	chatHandler := handler.NewChatHandler(authService, *chatService)

	// // InMemoryStateManager는 클라이언트와 채팅방의 상태를 메모리에 저장하고 관리합니다.
	// // 이 메서드는 단순하고 빠르며 테스트 목적으로 적합합니다.
	// // 전략 패턴을 사용하여 이 부분을 다른 상태 관리 구현체로 쉽게 교체할 수 있습니다.
	// stateManager := chat.NewInMemoryStateManager()
	// wsServer := chat.NewWebSocketServer(stateManager)
	// go wsServer.Run()
	// chat.InitializeWebSocketServer(ctx, wsServer, chatService)
	// chatHandler := handler.NewChatHandler(wsServer, stateManager, authService, *chatService)

	// RegisterChan middlewares
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
		userAPIGroup.GET("/:userID", userHandler.FindUserByID)
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

	upgrader := wschat.NewDefaultUpgrader()
	wsServerV2 := wschat.NewWSServer(upgrader, authService, *mediaService, *chatService)

	go wsServerV2.LoopOverClientMessages()

	chatAPIGroup := apiRouteGroup.Group("/chat")
	{
		chatAPIGroup.GET("/ws", wsServerV2.HandleConnections)
		chatAPIGroup.POST("/rooms", chatHandler.CreateRoom)
		chatAPIGroup.PUT("/rooms/:roomID/join", chatHandler.JoinChatRoom)
		chatAPIGroup.PUT("/rooms/:roomID/leave", chatHandler.LeaveChatRoom)
		chatAPIGroup.GET("/rooms/:roomID", chatHandler.FindRoomByID)
		chatAPIGroup.GET("/rooms", chatHandler.FindAllRooms)
		chatAPIGroup.GET("/rooms/:roomID/messages", chatHandler.FindMessagesByRoomID)
	}

	return e, nil
}
