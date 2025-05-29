package cognito

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
	"github.com/aws/aws-sdk-go-v2/config"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/victorsvart/gofin/internal/domain/apperror"
	"github.com/victorsvart/gofin/internal/domain/structs"
)

type cognitoActions struct {
	Client *cognito.Client
}

func NewCognitoActor() structs.Cognito {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		panic(err)
	}

	client := cognito.NewFromConfig(cfg)
	return &cognitoActions{Client: client}
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

func (actor cognitoActions) SignIn(ctx context.Context,
	userName string,
	password string,
) (*types.AuthenticationResultType, *apperror.AppError) {
	var authResult *types.AuthenticationResultType
	output, err := actor.Client.InitiateAuth(ctx, &cognito.InitiateAuthInput{
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

	return authResult, nil
}

func (actor cognitoActions) SignUp(ctx context.Context,
	userName string,
	password string,
	userEmail string,
	userFullName string,
) (bool, *apperror.AppError) {
	confirmed := false
	output, err := actor.Client.SignUp(ctx, &cognito.SignUpInput{
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
			return false, apperror.NewAppError(
				userName,
				apperror.AUTH,
				apperror.INVALID,
				errors.New(*invalidPassword.Message),
			)
		} else {
			return false, apperror.NewAppError(userName, apperror.AUTH, apperror.INTERNAL, err)
		}
	} else {
		confirmed = output.UserConfirmed
	}

	return confirmed, nil
}

func (actor cognitoActions) ConfirmEmail(ctx context.Context, userName, confirmationCode string) *apperror.AppError {
	_, err := actor.Client.ConfirmSignUp(ctx, &cognito.ConfirmSignUpInput{
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

func (actor cognitoActions) AdminCreateUser(
	ctx context.Context,
	userName string,
	userEmail string,
) *apperror.AppError {
	_, err := actor.Client.AdminCreateUser(ctx, &cognito.AdminCreateUserInput{
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

func (actor cognitoActions) AdminSetUserPassword(
	ctx context.Context,
	userName string,
	password string,
) *apperror.AppError {
	_, err := actor.Client.AdminSetUserPassword(ctx, &cognito.AdminSetUserPasswordInput{
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

func (actor cognitoActions) ConfirmForgotPassword(
	ctx context.Context,
	code string,
	userName string,
	password string,
) *apperror.AppError {
	_, err := actor.Client.ConfirmForgotPassword(ctx, &cognito.ConfirmForgotPasswordInput{
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

func (actor cognitoActions) DeleteUser(ctx context.Context, userAccessToken string) *apperror.AppError {
	_, err := actor.Client.DeleteUser(ctx, &cognito.DeleteUserInput{
		AccessToken: aws.String(userAccessToken),
	})
	if err != nil {
		return apperror.NewAppError("undefined", apperror.AUTH, apperror.INTERNAL, err)
	}

	return nil
}

func (actor cognitoActions) ForgotPassword(ctx context.Context,
	clientId string,
	userName string,
) (*types.CodeDeliveryDetailsType, *apperror.AppError) {
	output, err := actor.Client.ForgotPassword(ctx, &cognito.ForgotPasswordInput{
		ClientId: aws.String(clientId),
		Username: aws.String(userName),
	})

	if err != nil {
		return nil, apperror.NewAppError(userName, apperror.AUTH, apperror.INTERNAL, err)
	}
	return output.CodeDeliveryDetails, nil
}
