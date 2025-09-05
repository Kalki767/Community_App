package http

import (
	"auth/internal/delivery/http/handlers"
	"auth/internal/delivery/http/middleware"
	usecaseinterfaces "auth/internal/domain/contracts/usecase_interfaces"
	"auth/internal/services"

	"github.com/gin-gonic/gin"
)

// RouterConfig holds all the handlers required by the router.
type RouterConfig struct {
    UserHandler    *handlers.UserHandler
    SessionHandler *handlers.SessionHandler
    TokenService services.TokenService
    SessionUsecase usecaseinterfaces.SessionUsecaseInterface
}

// SetupRouter configures and returns the Gin router.
func SetupRouter(config *RouterConfig) *gin.Engine {
    router := gin.New()

    // Define API routes
    api := router.Group("/api/v1")
    {
        // Public routes
        public := api.Group("/auth")
        {
            public.POST("/register", config.UserHandler.Register)
            public.POST("/login", config.UserHandler.Login)
            public.POST("/refresh", config.SessionHandler.Refresh)
        }

        // Protected routes (will need an auth middleware)
        protected := api.Group("/user")
        protected.Use(middleware.AuthMiddleware(config.TokenService,config.SessionUsecase))
        {
            protected.GET("/me", config.UserHandler.GetMe)
            protected.GET("/is-verified", config.UserHandler.IsVerified)
        }

        // More protected routes
        sessionRoutes := api.Group("/sessions")
        sessionRoutes.Use(middleware.AuthMiddleware(config.TokenService, config.SessionUsecase))
        {
            sessionRoutes.GET("/me", config.SessionHandler.ListActiveSessions)
            sessionRoutes.GET("/get-session", config.SessionHandler.GetSession)
            sessionRoutes.DELETE("/logout", config.SessionHandler.Logout)
            sessionRoutes.DELETE("/all-except", config.SessionHandler.LogoutAllExcept)
        }
    }

    return router
}