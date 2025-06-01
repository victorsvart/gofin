package cognitousecase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	provider "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/victorsvart/gofin/internal/domain/apperror"
	"github.com/victorsvart/gofin/internal/domain/structs/cognito"
)

type cognitoUseCases struct {
	Client *provider.Client
}

func NewCognitoUseCases(client *provider.Client) cognito.CognitoUseCases {
	return &cognitoUseCases{Client: client}
}

func userPoolID() string {
	pid := os.Getenv("USER_POOL_ID")
	if pid == "" {
		log.Fatalln("user pool not set")
	}

	return pid
}

func clientID() string {
	cid := os.Getenv("CLIENT_ID")
	if cid == "" {
		log.Fatalln("client not set")
	}

	return cid
}

func clientSecret() string {
	cs := os.Getenv("CLIENT_SECRET")
	if cs == "" {
		log.Fatalln("secret not set")
	}

	return cs
}

func computeSecretHash(username string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret()))
	mac.Write([]byte(username + clientID()))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (actor cognitoUseCases) SignIn(ctx context.Context,
	userName string,
	password string,
) (*cognito.SignInResponse, *apperror.AppError) {
	var authResult *types.AuthenticationResultType
	output, err := actor.Client.InitiateAuth(ctx, &provider.InitiateAuthInput{
		AuthFlow:       "USER_PASSWORD_AUTH",
		ClientId:       aws.String(clientID()),
		AuthParameters: map[string]string{"USERNAME": userName, "PASSWORD": password},
	})

	if err != nil {
		var resetRequired *types.PasswordResetRequiredException
		if errors.As(err, &resetRequired) {
			return nil, apperror.NewAppError(
				userName,
				apperror.AUTH,
				apperror.INVALID,
				errors.New(*resetRequired.Message),
			)
		} else {
			return nil, apperror.NewAppError(userName, apperror.AUTH, apperror.INTERNAL, err)
		}
	} else {
		authResult = output.AuthenticationResult
	}

	response := cognito.MapToSignInResponse(authResult)
	return &response, nil
}

func (actor cognitoUseCases) SignUp(ctx context.Context,
	userName string,
	password string,
	userEmail string,
	userFullName string,
) (*cognito.SignUpResponse, *apperror.AppError) {
	var response cognito.SignUpResponse
	output, err := actor.Client.SignUp(ctx, &provider.SignUpInput{
		ClientId:   aws.String(clientID()),
		Password:   aws.String(password),
		Username:   aws.String(userName),
		SecretHash: aws.String(computeSecretHash(userName)),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: aws.String(userEmail)},
			{Name: aws.String("name"), Value: aws.String(userFullName)},
		},
	})

	if err != nil {
		var invalidPassword *types.InvalidPasswordException
		if errors.As(err, &invalidPassword) {
			return nil, apperror.NewAppError(
				userName,
				apperror.AUTH,
				apperror.INVALID,
				errors.New(*invalidPassword.Message),
			)
		} else {
			return nil, apperror.NewAppError(userName, apperror.AUTH, apperror.INTERNAL, err)
		}
	} else {
		response.EmailConfirmed = output.UserConfirmed
	}

	return &response, nil
}

func (actor cognitoUseCases) ConfirmEmail(ctx context.Context, userName, confirmationCode string) *apperror.AppError {
	_, err := actor.Client.ConfirmSignUp(ctx, &provider.ConfirmSignUpInput{
		ClientId:         aws.String(clientID()),
		ConfirmationCode: aws.String(confirmationCode),
		Username:         aws.String(userName),
		SecretHash:       aws.String(computeSecretHash(userName)),
	})

	if err != nil {
		var codeMismatch *types.CodeMismatchException
		if errors.As(err, &codeMismatch) {
			return apperror.NewAppError(userName, apperror.EMAIL_CONFIRM, apperror.MISMATCH, errors.New("confirmation code mismatch"))
		}

		var expiredCode *types.ExpiredCodeException
		if errors.As(err, &expiredCode) {
			return apperror.NewAppError(userName, apperror.EMAIL_CONFIRM, apperror.EXPIRED, errors.New("code is expired"))
		}

		log.Println(err)
		return apperror.NewAppError(userName, apperror.EMAIL_CONFIRM, apperror.INTERNAL, err)
	}

	return nil
}

func (actor cognitoUseCases) AdminCreateUser(
	ctx context.Context,
	userName string,
	userEmail string,
) *apperror.AppError {
	_, err := actor.Client.AdminCreateUser(ctx, &provider.AdminCreateUserInput{
		UserPoolId:     aws.String(userPoolID()),
		Username:       aws.String(userName),
		MessageAction:  types.MessageActionTypeSuppress,
		UserAttributes: []types.AttributeType{{Name: aws.String("email"), Value: aws.String(userEmail)}},
	})
	if err != nil {
		var userExists *types.UsernameExistsException

		if errors.As(err, &userExists) {
			return apperror.NewAppError(userName, apperror.AUTH, apperror.EXISTS, errors.New("user already exists"))
		} else {
			return apperror.NewAppError(userName, apperror.AUTH, apperror.INTERNAL, err)
		}
	}
	return nil
}

func (actor cognitoUseCases) AdminSetUserPassword(
	ctx context.Context,
	userName string,
	password string,
) *apperror.AppError {
	_, err := actor.Client.AdminSetUserPassword(ctx, &provider.AdminSetUserPasswordInput{
		Password:   aws.String(password),
		UserPoolId: aws.String(userPoolID()),
		Username:   aws.String(userName),
		Permanent:  true,
	})

	if err != nil {
		var invalidPassword *types.InvalidPasswordException
		if errors.As(err, &invalidPassword) {
			message := *invalidPassword.Message
			return apperror.NewAppError(password, apperror.AUTH, apperror.INTERNAL, errors.New(message))
		} else {
			return apperror.NewAppError(password, apperror.AUTH, apperror.INTERNAL, err)
		}
	}

	return nil
}

func (actor cognitoUseCases) ConfirmForgotPassword(
	ctx context.Context,
	code string,
	userName string,
	password string,
) *apperror.AppError {
	_, err := actor.Client.ConfirmForgotPassword(ctx, &provider.ConfirmForgotPasswordInput{
		ClientId:         aws.String(clientID()),
		ConfirmationCode: aws.String(code),
		Password:         aws.String(password),
		Username:         aws.String(userName),
	})

	if err != nil {
		var invalidPassword *types.InvalidPasswordException
		if errors.As(err, &invalidPassword) {
			return apperror.NewAppError(
				fmt.Sprintf("%v - %v", userName, code),
				apperror.AUTH,
				apperror.INVALID,
				errors.New(*invalidPassword.Message),
			)
		} else {
			return apperror.NewAppError(
				fmt.Sprintf("%v - %v", userName, code),
				apperror.AUTH,
				apperror.INTERNAL,
				err,
			)
		}
	}

	return nil
}

func (actor cognitoUseCases) DeleteUser(ctx context.Context, userAccessToken string) *apperror.AppError {
	_, err := actor.Client.DeleteUser(ctx, &provider.DeleteUserInput{
		AccessToken: aws.String(userAccessToken),
	})
	if err != nil {
		return apperror.NewAppError("undefined", apperror.AUTH, apperror.INTERNAL, err)
	}

	return nil
}

func (actor cognitoUseCases) ForgotPassword(ctx context.Context,
	clientId string,
	userName string,
) (*types.CodeDeliveryDetailsType, *apperror.AppError) {
	output, err := actor.Client.ForgotPassword(ctx, &provider.ForgotPasswordInput{
		ClientId: aws.String(clientId),
		Username: aws.String(userName),
	})

	if err != nil {
		return nil, apperror.NewAppError(userName, apperror.AUTH, apperror.INTERNAL, err)
	}
	return output.CodeDeliveryDetails, nil
}
