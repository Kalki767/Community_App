package dto

import "github.com/google/uuid"

type UserDto struct {
	ID          uuid.UUID
	FullName    string
	Email       string
	Username    string
	PhoneNumber string
	Country     string
	IsVerified  bool
}

type RegisterUser struct {
	FullName    string
	Email       string
	Username    string
	Password string
	PhoneNumber string
	Country     string
}

type LoginRequest struct {
	Identification string
	Password string
}