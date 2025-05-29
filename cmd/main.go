package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/victorsvart/gofin/internal/handler"
	"github.com/victorsvart/gofin/internal/infra/cognito"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	cognito := cognito.NewCognitoActor()
	authHandler := handler.NewAuthHandler(cognito)

	server := gin.Default()
	api := server.Group("/api/auth")
	{
		api.POST("signin", authHandler.SignIn)
		api.POST("signup", authHandler.SignUp)
		api.POST("confirmemail", authHandler.ConfirmEmail)
	}

	if err := server.Run(); err != nil {
		panic(err)
	}
}
