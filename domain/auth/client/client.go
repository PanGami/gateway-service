package client

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/pangami/gateway-service/domain/auth"

	rest "github.com/pangami/gateway-service/util/client"
)

var (
	loadEnvOnce sync.Once
	envLoadErr  error
)

func LoadEnv() error {
	loadEnvOnce.Do(func() {
		envLoadErr = godotenv.Load()
	})
	if envLoadErr != nil {
		return fmt.Errorf("error loading .env file: %v", envLoadErr)
	}
	return nil
}

type InitClient struct {
	Client  *rest.RestClient
	BaseURL string
	ApiKey  string
}

func InitializeClient() (InitClient, error) {
	baseURL := os.Getenv("AUTH_SERVICE_BASE_URL")
	apiKey := os.Getenv("AUTH_SERVICE_API_KEY")
	if baseURL == "" || apiKey == "" {
		return InitClient{}, errors.New("environment variables AUTH_SERVICE_BASE_URL or AUTH_SERVICE_API_KEY not set")
	}

	client := rest.NewRestClient(apiKey)
	return InitClient{
		Client:  client,
		BaseURL: baseURL,
		ApiKey:  apiKey,
	}, nil
}

func performRequest(ctx context.Context, endpoint string, payload map[string]interface{}, segment, subSegment string) (map[string]interface{}, error) {
	client, err := InitializeClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client: %v", err)
	}

	url := fmt.Sprintf("%s%s", client.BaseURL, endpoint)

	resClient, err := client.Client.CallAPI("POST", url, payload)
	if err != nil {
		return nil, err
	}

	return resClient, nil
}

func Login(ctx context.Context, request *auth.LoginRequest) (map[string]interface{}, error) {
	if err := LoadEnv(); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %v", err)
	}

	payload := map[string]interface{}{
		"username": request.Username,
		"password": request.Password,
		"fcmToken": request.FcmToken,
	}
	return performRequest(ctx, "login", payload, "auth", "login")
}

func Logout(ctx context.Context, request *auth.LogoutRequest) (map[string]interface{}, error) {
	if err := LoadEnv(); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %v", err)
	}

	payload := map[string]interface{}{
		"userId": nil,
		"token":  request.Token,
	}
	return performRequest(ctx, "logout", payload, "auth", "logout")
}

func ValidateToken(ctx context.Context, request string) (map[string]interface{}, error) {
	if err := LoadEnv(); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %v", err)
	}

	payload := map[string]interface{}{
		"token": request,
	}
	return performRequest(ctx, "validate_token", payload, "auth", "validate_token")
}
