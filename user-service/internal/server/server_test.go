package server

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"grpc-demo/user-service/internal/storage"

	pb "grpc-demo/pb"
)

// newTestClient поднимает настоящий gRPC-сервер в памяти (bufconn)
// и возвращает клиента, так что тесты проходят весь путь запроса.
func newTestClient(t *testing.T) pb.UserServiceClient {
	t.Helper()

	lis := bufconn.Listen(1024 * 1024)
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, NewUserServer(storage.NewInMemoryStorage()))

	go func() { _ = grpcServer.Serve(lis) }()
	t.Cleanup(grpcServer.Stop)

	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	return pb.NewUserServiceClient(conn)
}

func TestCreateAndGetUser(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	created, err := c.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "Anna",
		Email: "anna@example.com",
	})
	require.NoError(t, err)
	require.NotEmpty(t, created.GetId())

	got, err := c.GetUser(ctx, &pb.GetUserRequest{Id: created.GetId()})
	require.NoError(t, err)
	require.Equal(t, "Anna", got.GetName())
	require.Equal(t, "anna@example.com", got.GetEmail())
}

func TestGetUser_Errors(t *testing.T) {
	c := newTestClient(t)

	tests := []struct {
		name     string
		id       string
		wantCode codes.Code
	}{
		{name: "несуществующий пользователь", id: "unknown-id", wantCode: codes.NotFound},
		{name: "пустой id", id: "", wantCode: codes.InvalidArgument},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := c.GetUser(context.Background(), &pb.GetUserRequest{Id: tt.id})
			require.Error(t, err)
			require.Equal(t, tt.wantCode, status.Code(err))
		})
	}
}

func TestCreateUser_Validation(t *testing.T) {
	c := newTestClient(t)

	tests := []struct {
		name     string
		req      *pb.CreateUserRequest
		wantCode codes.Code
	}{
		{name: "корректные данные", req: &pb.CreateUserRequest{Name: "Anna", Email: "anna@example.com"}, wantCode: codes.OK},
		{name: "пустое имя", req: &pb.CreateUserRequest{Name: "  ", Email: "anna@example.com"}, wantCode: codes.InvalidArgument},
		{name: "пустой email", req: &pb.CreateUserRequest{Name: "Anna", Email: ""}, wantCode: codes.InvalidArgument},
		{name: "некорректный email", req: &pb.CreateUserRequest{Name: "Anna", Email: "not-an-email"}, wantCode: codes.InvalidArgument},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := c.CreateUser(context.Background(), tt.req)
			require.Equal(t, tt.wantCode, status.Code(err))
		})
	}
}
