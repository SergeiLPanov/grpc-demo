package client

import (
	"context"
	"fmt"
	"time"

	pb "grpc-demo/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// callTimeout ограничивает время одного вызова user-service.
const callTimeout = 3 * time.Second

type UserClient struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
}

func NewUserClient(addr string) (*UserClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать клиент user-service: %w", err)
	}

	return &UserClient{
		conn:   conn,
		client: pb.NewUserServiceClient(conn),
	}, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}

// CreateUser создаёт пользователя в user-service и возвращает его id.
func (c *UserClient) CreateUser(ctx context.Context, name, email string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()

	resp, err := c.client.CreateUser(ctx, &pb.CreateUserRequest{Name: name, Email: email})
	if err != nil {
		return "", fmt.Errorf("ошибка вызова CreateUser: %w", err)
	}
	return resp.GetId(), nil
}

// UserExists проверяет, есть ли пользователь. Статус NotFound — не ошибка,
// а ответ «пользователя нет»; остальные ошибки возвращаются вызывающему.
func (c *UserClient) UserExists(ctx context.Context, userID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()

	_, err := c.client.GetUser(ctx, &pb.GetUserRequest{Id: userID})
	switch status.Code(err) {
	case codes.OK:
		return true, nil
	case codes.NotFound:
		return false, nil
	default:
		return false, fmt.Errorf("ошибка вызова GetUser: %w", err)
	}
}
