package middleware

import (
	usecaseinterfaces "auth/internal/domain/contracts/usecase_interfaces"
	"auth/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenSerivce services.TokenService, sessionUsecase usecaseinterfaces.SessionUsecaseInterface) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == ""{
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer"{
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization header format"})
			return
		}

		accessToken := parts[1]
		claims, err := tokenSerivce.ParseAccessToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token", "error": err.Error()})
			return
		}

		active, err := sessionUsecase.IsSessionActive(claims.SessionID)
		if err != nil || !active {
    		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Session expired or revoked"})
    		return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("session_id", claims.SessionID)

		ctx.Next()
	}
}