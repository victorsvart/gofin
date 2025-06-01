package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type cognitoClient struct {
	Client *cognito.Client
}

func NewCognitoClient() *cognito.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		panic(err)
	}

	client := cognito.NewFromConfig(cfg)
	return client
}
