package request

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OK(c *gin.Context, a any) {
	c.JSON(http.StatusBadRequest, a)
}

func BadRequest(c *gin.Context, a any) {
	c.JSON(http.StatusBadRequest, gin.H{"error": a})
}

func Unauthorized(c *gin.Context, a any) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": a})
}

func InternalServerError(c *gin.Context, a any) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": a})
}
