package handlers

import (
	"errors"
	"net/http"

	"auth/internal/delivery/http/dto"
	usecaseinterfaces "auth/internal/domain/contracts/usecase_interfaces"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SessionHandler defines the HTTP handlers for sessions.
type SessionHandler struct {
	usecase usecaseinterfaces.SessionUsecaseInterface
}

// NewSessionHandler creates a new instance of SessionHandler.
func NewSessionHandler(usecase usecaseinterfaces.SessionUsecaseInterface) *SessionHandler {
	return &SessionHandler{usecase: usecase}
}

// ListActiveSessions godoc
// @Summary      Get all active sessions for a user
// @Description  Retrieves a list of all active sessions for the authenticated user.
// @Tags         sessions
// @Produce      json
// @Success      200  {array}   dto.SessionResponseDTO
// @Failure      401  {object}  dto.MessageResponse
// @Failure      500  {object}  dto.MessageResponse
// @Security     Bearer
// @Router       /sessions/me [get]
func (h *SessionHandler) ListActiveSessions(ctx *gin.Context) {
	userID, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "User ID not found in context"})
		return
	}

	parsedID, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "User ID in context is not a valid UUID"})
		return
	}

	sessions, err := h.usecase.ListActiveSessions(parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found", "error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve sessions", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, sessions)
}

// GetSession godoc
// @Summary      Get a specific session
// @Description  Retrieves a specific session by its ID for the authenticated user.
// @Tags         sessions
// @Produce      json
// @Success      200  {object}  dto.SessionResponseDTO
// @Failure      401  {object}  dto.MessageResponse
// @Failure      404  {object}  dto.MessageResponse
// @Failure      500  {object}  dto.MessageResponse
// @Security     Bearer
// @Router       /sessions/get-session [get]
func (h *SessionHandler) GetSession(ctx *gin.Context) {
	sessionID, ok := ctx.Get("session_id")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Session ID not found in context"})
		return
	}

	parsedID, ok := sessionID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Session ID in context is not a valid UUID"})
		return
	}

	sessionDTO, err := h.usecase.GetSession(parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Session not found", "error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve session", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, sessionDTO)
}

// Logout godoc
// @Summary      Logout a specific session
// @Description  Logs out the authenticated session by its ID.
// @Tags         sessions
// @Produce      json
// @Success      200  {object}  dto.MessageResponse
// @Failure      401  {object}  dto.MessageResponse
// @Failure      404  {object}  dto.MessageResponse
// @Failure      500  {object}  dto.MessageResponse
// @Security     Bearer
// @Router       /sessions/logout [delete]
func (h *SessionHandler) Logout(ctx *gin.Context) {
	sessionID, ok := ctx.Get("session_id")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Session ID not found in context"})
		return
	}

	parsedID, ok := sessionID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Session ID in context is not a valid UUID"})
		return
	}

	err := h.usecase.Logout(parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Session not found", "error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to log out session", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// LogoutAllExcept godoc
// @Summary      Logout all sessions except the current one
// @Description  Logs out all other active sessions for the authenticated user.
// @Tags         sessions
// @Produce      json
// @Success      200  {object}  dto.MessageResponse
// @Failure      401  {object}  dto.MessageResponse
// @Failure      404  {object}  dto.MessageResponse
// @Failure      500  {object}  dto.MessageResponse
// @Security     Bearer
// @Router       /sessions/all-except [delete]
func (h *SessionHandler) LogoutAllExcept(ctx *gin.Context) {
	userID, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "User ID not found in context"})
		return
	}

	parsedUserID, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "User ID in context is not a valid UUID"})
		return
	}

	sessionID, ok := ctx.Get("session_id")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Session ID not found in context"})
		return
	}

	parsedSessionID, ok := sessionID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Session ID in context is not a valid UUID"})
		return
	}

	err := h.usecase.LogoutAllExcept(parsedUserID, parsedSessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "User or session not found", "error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to log out sessions", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out from other devices except this one"})
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Refreshes an expired access token using a valid refresh token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RefreshRequest  true  "Refresh token"
// @Success      200      {object}  dto.RefreshResponse
// @Failure      400      {object}  dto.MessageResponse
// @Failure      401      {object}  dto.MessageResponse
// @Failure      404      {object}  dto.MessageResponse
// @Failure      500      {object}  dto.MessageResponse
// @Router       /auth/refresh [post]
func (h *SessionHandler) Refresh(ctx *gin.Context) {
	var req dto.RefreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
		return
	}

	accessToken, err := h.usecase.Refresh(req.RefreshToken)
	if err != nil {
		switch err.Error() {
		case "session expired or revoked":
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Session expired or revoked", "error": err.Error()})
		case "refresh token mismatch":
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid refresh token", "error": err.Error()})
		default:
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"message": "Session not found", "error": err.Error()})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to refresh token", "error": err.Error()})
			}
		}
		return
	}

	ctx.JSON(http.StatusOK, dto.RefreshResponse{AccessToken: accessToken})
}
