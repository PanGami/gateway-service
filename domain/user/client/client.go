package client

import (
	"context"
	"log"
	"os"

	pb "github.com/pangami/gateway-service/proto/user"
	Grpc "github.com/pangami/gateway-service/util/client"
)

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

func UserUpdate(ctx context.Context, req *pb.CreateUserRequest) (*pb.NoResponse, error) {
	conn := Grpc.Dial(os.Getenv("USER_SERVICE_GRPC"))
	defer conn.Close()

	client := pb.NewUserClient(conn)

	res, err := client.UpdateUser(ctx, req)
	if err != nil {
		log.Printf("Create User gRPC error: %v\n", err)
		return nil, err
	}

	return res, nil
}

func UserDelete(ctx context.Context, req *pb.DetailUserRequest) (*pb.NoResponse, error) {
	conn := Grpc.Dial(os.Getenv("USER_SERVICE_GRPC"))
	defer conn.Close()

	client := pb.NewUserClient(conn)

	res, err := client.DeleteUser(ctx, req)
	if err != nil {
		log.Printf("Create User gRPC error: %v\n", err)
		return nil, err
	}

	return res, nil
}

func UserDetail(ctx context.Context, req *pb.DetailUserRequest) (*pb.DetailUserResponse, error) {
	conn := Grpc.Dial(os.Getenv("USER_SERVICE_GRPC"))
	defer conn.Close()

	client := pb.NewUserClient(conn)

	res, err := client.DetailUser(ctx, req)
	if err != nil {
		log.Printf("Create User gRPC error: %v\n", err)
		return nil, err
	}

	return res, nil
}

func UserList(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	conn := Grpc.Dial(os.Getenv("USER_SERVICE_GRPC"))
	defer conn.Close()

	client := pb.NewUserClient(conn)

	res, err := client.ListUsers(ctx, req)
	if err != nil {
		log.Printf("List Users gRPC error: %v\n", err)
		return nil, err
	}

	return res, nil
}

func UserGetActivity(ctx context.Context, req *pb.DetailUserRequest) (*pb.UserActivitiesResponse, error) {
	conn := Grpc.Dial(os.Getenv("USER_SERVICE_GRPC"))
	defer conn.Close()

	client := pb.NewUserClient(conn)

	res, err := client.GetUserActivities(ctx, req)
	if err != nil {
		log.Printf("Get Activity gRPC error: %v\n", err)
		return nil, err
	}

	return res, nil
}
