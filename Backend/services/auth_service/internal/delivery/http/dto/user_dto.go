package dto

import "github.com/google/uuid"

type UserDto struct {
	ID          uuid.UUID `json:"id"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	PhoneNumber string    `json:"phone_number"`
	Country     string    `json:"country"`
	IsVerified  bool      `json:"is_verified"`
}

type RegisterUser struct {
	FullName    string `json:"full_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	PhoneNumber string `json:"phone_number"`
	Country     string `json:"country"`
}

type LoginRequest struct {
	Identification string `json:"identification" binding:"required"`
	Password       string `json:"password" binding:"required"`
}
