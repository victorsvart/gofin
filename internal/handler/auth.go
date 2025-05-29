package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorsvart/gofin/internal/domain/structs"
	"github.com/victorsvart/gofin/internal/infra/request"
)

type AuthHandler struct {
	Cognito structs.Cognito
}

func NewAuthHandler(cognito structs.Cognito) *AuthHandler {
	return &AuthHandler{Cognito: cognito}
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var req structs.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		request.BadRequest(c, err)
		return
	}

	authResult, appErr := h.Cognito.SignIn(c.Request.Context(), req.Username, req.Password)
	if appErr != nil {
		request.Unauthorized(c, appErr)
		return
	}

	request.OK(c, *authResult)
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req structs.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		request.BadRequest(c, err)
		return
	}

	confirmed, appErr := h.Cognito.SignUp(c.Request.Context(), req.Username, req.Password, req.EmailAddress, req.Name)
	if appErr != nil {
		request.InternalServerError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"confirmed": confirmed})
}

func (h *AuthHandler) ConfirmEmail(c *gin.Context) {
	var req structs.ConfirmEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		request.BadRequest(c, err)
		return
	}

	if err := h.Cognito.ConfirmEmail(c.Request.Context(), req.Username, req.ConfirmationCode); err != nil {
		request.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
