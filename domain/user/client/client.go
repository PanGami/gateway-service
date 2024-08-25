package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	user "github.com/pangami/gateway-service/domain/user"
	rest "github.com/pangami/gateway-service/util/client"

	pb "github.com/pangami/gateway-service/proto/user"
	Grpc "github.com/pangami/gateway-service/util/client"
)

type InitClient struct {
	Client  *rest.RestClient
	BaseURL string
	ApiKey  string
}

var loadEnvOnce sync.Once
var envLoadErr error

func LoadEnv() error {
	loadEnvOnce.Do(func() {
		envLoadErr = godotenv.Load()
	})
	if envLoadErr != nil {
		return fmt.Errorf("error loading .env file: %v", envLoadErr)
	}
	return nil
}

func InitializeClient() (InitClient, error) {
	baseURL := os.Getenv("USER_SERVICE_BASE_URL")
	apiKey := os.Getenv("USER_SERVICE_API_KEY")
	if baseURL == "" || apiKey == "" {
		return InitClient{}, fmt.Errorf("environment variables USER_SERVICE_BASE_URL or USER_SERVICE_API_KEY not set")
	}
	client := rest.NewRestClient(apiKey)
	return InitClient{
		Client:  client,
		BaseURL: baseURL,
		ApiKey:  apiKey,
	}, nil
}

func performRequest(ctx context.Context, client InitClient, endpoint string, payload map[string]interface{}, method string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s%s", client.BaseURL, endpoint)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
	}()

	res, err := client.Client.CallAPI(method, url, payload)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func UserList(ctx context.Context, filter *user.ListFilter) (map[string]interface{}, error) {
	client, err := InitializeClient()
	if err != nil {
		log.Fatal(err)
	}

	endpoint := "list"
	payload := map[string]interface{}{}
	return performRequest(ctx, client, endpoint, payload, "GET")
}

func UserCreate(ctx context.Context, req *pb.CreateUserRequest) (*pb.NoResponse, error) {
	conn := Grpc.Dial(os.Getenv("USER_SERVICE_GRPC"))
	defer conn.Close()

	client := pb.NewUserClient(conn)

	res, err := client.CreateUser(ctx, req)
	if err != nil {
		log.Printf("Create User gRPC error: %v\n", err)
		return nil, err
	}

	return res, nil
}

func UserUpdate(ctx context.Context, user *user.User) (map[string]interface{}, error) {
	client, err := InitializeClient()
	if err != nil {
		log.Fatal(err)
	}

	endpoint := fmt.Sprintf("update?id=%d", user.ID)
	payload := map[string]interface{}{
		"username": user.Username,
		"password": user.Password,
		"fullname": user.FullName,
	}
	return performRequest(ctx, client, endpoint, payload, "PUT")
}

func UserDelete(ctx context.Context, userID string) (map[string]interface{}, error) {
	client, err := InitializeClient()
	if err != nil {
		log.Fatal(err)
	}

	endpoint := fmt.Sprintf("delete?id=%s", userID)
	payload := map[string]interface{}{}
	return performRequest(ctx, client, endpoint, payload, "DELETE")
}

func UserDetail(ctx context.Context, userID string) (map[string]interface{}, error) {
	client, err := InitializeClient()
	if err != nil {
		log.Fatal(err)
	}

	endpoint := fmt.Sprintf("detail?id=%s", userID)
	payload := map[string]interface{}{}
	return performRequest(ctx, client, endpoint, payload, "GET")
}
