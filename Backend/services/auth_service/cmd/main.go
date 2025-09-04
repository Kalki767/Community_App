package main

import (
	"auth/internal/delivery/http"
	"log"
	"os"
	"time"
	// "github.com/joho/godotenv"

	handlers "auth/internal/delivery/http/handlers"
	repository "auth/internal/repository"
	services "auth/internal/services"
	usecase "auth/internal/usecase"

	"github.com/gin-gonic/gin"

	// swagger
	_ "auth/cmd/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	db "auth/internal/infrastructure/db"
	cache "auth/internal/infrastructure/cache"

)

// @title           Authentication Service API
// @version         1.0
// @description     This is the API for the Community App's authentication service.
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apiKey  Bearer
// @in header
// @name Authorization
func main() {
// 	if _, err := os.Stat(".env"); err == nil {
//     if err := godotenv.Load("../.env"); err != nil {
//         log.Println("Failed to load .env file:", err)
//     }
// }

	// if err := godotenv.Load("../.env"); err != nil {
    //     log.Println("No .env file found, relying on system environment variables")
    // }
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)

	// --- 1. Database Connection and Migration ---
	database := db.NewPostgresDB()
	db.RunMigrations(database)

	redis := cache.NewRedisClient()
	log.Println(redis)

	// --- 2. Component Initialization (from bottom-up) ---
	// Repositories
	userRepo := repository.NewUserRepo(database)
	sessionRepo := repository.NewSessionRepository(database)

	// Services
	// TODO: Replace these hardcoded secrets with environment variables in a production setup.
	accessSecret := os.Getenv("ACCESS_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")
	tokenService := services.NewTokenService(accessSecret, refreshSecret, 15*time.Minute, 30*24*time.Hour)

	// Use Cases
	userUsecase := usecase.NewUserUsecase(userRepo, sessionRepo, tokenService)
	sessionUsecase := usecase.NewSessionUsecase(sessionRepo, tokenService)

	// Handlers
	userHandler := handlers.NewUserHandler(userUsecase)
	sessionHandler := handlers.NewSessionHandler(sessionUsecase)

	// --- 3. Route Configuration ---
	routerConfig := &http.RouterConfig{
		UserHandler:    userHandler,
		SessionHandler: sessionHandler,
		TokenService:   tokenService,
		SessionUsecase: sessionUsecase,
	}
	router := http.SetupRouter(routerConfig)

	// Add the Swagger endpoint
	docsUrl := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The URL to your swagger.json
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, docsUrl))

	// --- 4. Server Startup ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
