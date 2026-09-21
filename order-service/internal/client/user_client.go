package client

import (
	"context"
	"fmt"

	pb "grpc-demo/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
}

func NewUserClient(addr string) (*UserClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к user-service: %w", err)
	}

	return &UserClient{
		conn:   conn,
		client: pb.NewUserServiceClient(conn),
	}, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}

func (c *UserClient) UserExists(ctx context.Context, userID string) (bool, error) {
	resp, err := c.client.GetUser(ctx, &pb.GetUserRequest{Id: userID})
	if err != nil {
		return false, fmt.Errorf("ошибка вызова GetUser: %w", err)
	}
	return resp.Found, nil
}
