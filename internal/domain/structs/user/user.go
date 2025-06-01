package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/victorsvart/gofin/internal/domain/apperror"
)

type User struct {
	ID           uuid.UUID
	Name         string `json:"name" binding:"required"`
	Username     string `json:"username" binding:"required"`
	EmailAddress string `json:"emailAddress" binding:"required,email"`
}

type CreateUserResponse struct {
	ID       uuid.UUID `json:"sub"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
}

type CreateUserInput struct {
	ID           uuid.UUID `json:"sub"`
	Name         string    `json:"name"`
	Username     string    `json:"username"`
	EmailAddress string    `json:"emailAddress" binding:"required,email"`
}

type UserRepository interface {
	InsertUser(context.Context, CreateUserInput) (*CreateUserResponse, *apperror.AppError)
}
