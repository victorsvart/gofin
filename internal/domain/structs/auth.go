package structs

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

type ConfirmEmailRequest struct {
	Username         string `json:"username" binding:"required"`
	ConfirmationCode string `json:"code" binding:"required"`
}

type Cognito interface {
	SignIn(ctx context.Context, userName string, password string) (*types.AuthenticationResultType, *apperror.AppError)
	SignUp(ctx context.Context, userName string, password string, userEmail string, userFullName string) (bool, *apperror.AppError)
	ConfirmEmail(ctx context.Context, userName, confirmationCode string) *apperror.AppError
	AdminCreateUser(ctx context.Context, userName string, userEmail string) *apperror.AppError
	AdminSetUserPassword(ctx context.Context, userName string, password string) *apperror.AppError
	ConfirmForgotPassword(ctx context.Context, code string, userName string, password string) *apperror.AppError
	DeleteUser(ctx context.Context, userAccessToken string) *apperror.AppError
	ForgotPassword(ctx context.Context, clientId string, userName string) (*types.CodeDeliveryDetailsType, *apperror.AppError)
}
