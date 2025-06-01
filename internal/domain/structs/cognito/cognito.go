package cognito

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/victorsvart/gofin/internal/domain/apperror"
)

type SignInRequest struct {
	Username     string `json:"username" binding:"required"`
	EmailAddress string `json:"emailAddress" binding:"required,email"`
	Password     string `json:"password" binding:"required"`
}

type SignUpRequest struct {
	Name         string `json:"name" binding:"required"`
	Username     string `json:"username" binding:"required"`
	EmailAddress string `json:"emailAddress" binding:"required,email"`
	Password     string `json:"password" binding:"required"`
}

type SignInResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int32  `json:"expiresIn"`
	TokenType    string `json:"tokenType"`
}

type SignUpResponse struct {
	EmailConfirmed bool `json:"emailConfirmed"`
}

func MapToSignInResponse(authResult *types.AuthenticationResultType) SignInResponse {
	response := SignInResponse{
		ExpiresIn: authResult.ExpiresIn,
	}

	if authResult.AccessToken != nil {
		response.AccessToken = *authResult.AccessToken
	}
	if authResult.RefreshToken != nil {
		response.RefreshToken = *authResult.RefreshToken
	}
	if authResult.TokenType != nil {
		response.TokenType = *authResult.TokenType
	}

	return response
}

func MapToSignUpResponse(emailConfirmed bool) SignUpResponse {
	return SignUpResponse{EmailConfirmed: emailConfirmed}
}

type ConfirmEmailRequest struct {
	Username         string `json:"username" binding:"required"`
	ConfirmationCode string `json:"code" binding:"required"`
}

type CognitoUseCases interface {
	SignIn(ctx context.Context, userName string, password string) (*SignInResponse, *apperror.AppError)
	SignUp(ctx context.Context, userName string, password string, userEmail string, userFullName string) (*SignUpResponse, *apperror.AppError)
	ConfirmEmail(ctx context.Context, userName, confirmationCode string) *apperror.AppError
	AdminCreateUser(ctx context.Context, userName string, userEmail string) *apperror.AppError
	AdminSetUserPassword(ctx context.Context, userName string, password string) *apperror.AppError
	ConfirmForgotPassword(ctx context.Context, code string, userName string, password string) *apperror.AppError
	DeleteUser(ctx context.Context, userAccessToken string) *apperror.AppError
	ForgotPassword(ctx context.Context, clientId string, userName string) (*types.CodeDeliveryDetailsType, *apperror.AppError)
}
