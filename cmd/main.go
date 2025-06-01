package main

import (
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/victorsvart/gofin/docs"
	"github.com/victorsvart/gofin/internal/infra/gin/routing"
)

// @title           GoFin
// @version         1.0
// @description     Api for GoFin App.
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	server := routing.SetRoutes()
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := server.Run(); err != nil {
		panic(err)
	}
}
