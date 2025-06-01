package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/victorsvart/gofin/internal/adapters/handlers"
	"github.com/victorsvart/gofin/internal/infra/aws"
	"github.com/victorsvart/gofin/internal/usecases/cognitousecase"
)

var server = gin.Default()
var v1 *gin.RouterGroup = server.Group("/api/v1")

func SetRoutes() *gin.Engine {
	routeAuth()
	return server
}

func routeAuth() {
	cc := aws.NewCognitoClient()
	u := cognitousecase.NewCognitoUseCases(cc)
	handler := handlers.NewAuthHandler(u)
	r := v1.Group("/auth")
	{
		r.POST("signin", handler.SignIn)
		r.POST("signup", handler.SignUp)
		r.POST("confirmemail", handler.ConfirmEmail)
	}
}
