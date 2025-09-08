package handlers

import (
	dto "auth/internal/delivery/http/dto"
	usecaseinterfaces "auth/internal/domain/contracts/usecase_interfaces"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserHandler defines the HTTP handlers for user-related actions.
type UserHandler struct {
	userusecase usecaseinterfaces.UserUsecaseInterface
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(userusecase usecaseinterfaces.UserUsecaseInterface) *UserHandler {
	return &UserHandler{userusecase: userusecase}
}

// Register godoc
// @Summary      Register a new user
// @Description  Registers a new user with email, password, and other details.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      dto.RegisterUser  true  "User registration data"
// @Success      201  {object}  dto.UserDto
// @Failure      400  {object}  dto.MessageResponse
// @Failure      500  {object}  dto.MessageResponse
// @Router       /auth/register [post]
func (handler *UserHandler) Register(ctx *gin.Context) {
	var userdto dto.RegisterUser
	if err := ctx.ShouldBindJSON(&userdto); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request format", "error": err.Error()})
		return
	}

	user, err := handler.userusecase.Register(&userdto)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Cannot register user", "error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, gin.H{"message": "Successfully registered", "data": user})
}

// Login godoc
// @Summary      Login a user
// @Description  Authenticates a user and returns access and refresh tokens.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      dto.LoginRequest  true  "User login credentials"
// @Success      200  {object}  dto.UserDto
// @Failure      400  {object}  dto.MessageResponse
// @Failure      401  {object}  dto.MessageResponse
// @Failure      404  {object}  dto.MessageResponse
// @Router       /auth/login [post]
func (handler *UserHandler) Login(ctx *gin.Context) {
	var request dto.LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request format", "error": err.Error()})
		return
	}

	userdto, refreshtoken, accesstoken, err := handler.userusecase.Login(request.Identification, request.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "user not found" {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found", "error": err.Error()})
			return
		}
		if err.Error() == "invalid credentials" {
			ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials", "error": err.Error()})
			return
		}
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Cannot login to the system", "error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"message":      "Successfully logged in",
		"user":         userdto,
		"refreshtoken": refreshtoken,
		"accesstoken":  accesstoken,
	})
}

// GetMe godoc
// @Summary      Get authenticated user's profile
// @Description  Retrieves the profile of the user authenticated by the JWT token.
// @Tags         user
// @Produce      json
// @Success      200  {object}  dto.UserDto
// @Failure      401  {object}  dto.MessageResponse
// @Security     Bearer
// @Router       /user/me [get]
func (handler *UserHandler) GetMe(ctx *gin.Context) {
	userId, ok := ctx.Get("user_id")
	if !ok {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "User ID not found in context"})
		return
	}

	parsedId, ok := userId.(uuid.UUID)
	if !ok {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "User ID in context is not a valid UUID"})
		return
	}

	userdto, err := handler.userusecase.GetUserProfile(parsedId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found", "error": err.Error()})
			return
		}
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Cannot retrieve user profile", "error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully retrieved user", "data": userdto})
}

// IsVerified godoc
// @Summary      Check if authenticated user is verified
// @Description  Checks the verification status of the user authenticated by the JWT token.
// @Tags         user
// @Produce      json
// @Success      200  {object}  dto.UserDto
// @Failure      401  {object}  dto.MessageResponse
// @Failure      404  {object}  dto.MessageResponse
// @Security     Bearer
// @Router       /user/is-verified [get]
func (handler *UserHandler) IsVerified(ctx *gin.Context) {
	userId, ok := ctx.Get("user_id")
	if !ok {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "User ID not found in context"})
		return
	}

	parsedId, ok := userId.(uuid.UUID)
	if !ok {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "User ID in context is not a valid UUID"})
		return
	}

	isverified, err := handler.userusecase.IsVerifiedUser(parsedId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found", "error": err.Error()})
			return
		}
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Cannot check user verification status", "error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully retrieved user verification status", "data": isverified})
}
