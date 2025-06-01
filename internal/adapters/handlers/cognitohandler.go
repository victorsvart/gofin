package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorsvart/gofin/internal/domain/request"
	"github.com/victorsvart/gofin/internal/domain/structs/cognito"
)

// AuthHandler handles authentication routes.
type AuthHandler struct {
	CognitoUseCases cognito.CognitoUseCases
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(cognito cognito.CognitoUseCases) *AuthHandler {
	return &AuthHandler{CognitoUseCases: cognito}
}

// SignIn godoc
// @Summary Sign in a user
// @Description Authenticate a user and return access and refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param SignInRequest body cognito.SignInRequest true "Sign In Request"
// @Success 200 {object} cognito.SignInResponse
// @Failure 400 {object} apperror.AppError
// @Failure 401 {object} apperror.AppError
// @Router /auth/signin [post]
func (h *AuthHandler) SignIn(c *gin.Context) {
	var req cognito.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		request.BadRequest(c, err)
		return
	}

	authResult, appErr := h.CognitoUseCases.SignIn(c.Request.Context(), req.Username, req.Password)
	if appErr != nil {
		request.Unauthorized(c, appErr)
		return
	}

	request.OK(c, *authResult)
}

// SignUp godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param SignUpRequest body cognito.SignUpRequest true "Sign Up Request"
// @Success 200 {object} cognito.SignUpResponse
// @Failure 400 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /auth/signup [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req cognito.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		request.BadRequest(c, err)
		return
	}

	response, appErr := h.CognitoUseCases.SignUp(c.Request.Context(), req.Username, req.Password, req.EmailAddress, req.Name)
	if appErr != nil {
		request.InternalServerError(c, appErr)
		return
	}

	request.OK(c, response)
}

// ConfirmEmail godoc
// @Summary Confirm user's email
// @Description Confirm user account using a confirmation code
// @Tags Auth
// @Accept json
// @Produce json
// @Param ConfirmEmailRequest body cognito.ConfirmEmailRequest true "Email Confirmation Request"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /auth/confirmemail [post]
func (h *AuthHandler) ConfirmEmail(c *gin.Context) {
	var req cognito.ConfirmEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		request.BadRequest(c, err)
		return
	}

	if err := h.CognitoUseCases.ConfirmEmail(c.Request.Context(), req.Username, req.ConfirmationCode); err != nil {
		request.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
