package user

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID
	Name         string `json:"name" binding:"required"`
	Username     string `json:"username" binding:"required"`
	EmailAddress string `json:"emailAddress" binding:"required,email"`
	Password     string `json:"password"`
}
