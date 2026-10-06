package server

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"grpc-demo/user-service/internal/storage"

	pb "grpc-demo/pb"
)

// UserServer реализует сгенерированный интерфейс pb.UserServiceServer.
type UserServer struct {
	pb.UnimplementedUserServiceServer // обязательное встраивание для совместимости
	store                             storage.UserStorage
}

func NewUserServer(store storage.UserStorage) *UserServer {
	return &UserServer{store: store}
}

// GetUser возвращает пользователя по id.
// Ошибки отдаются как gRPC-статусы:
//   - InvalidArgument — пустой id;
//   - NotFound — пользователя нет;
//   - Internal — любая другая ошибка хранилища.
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	user, err := s.store.Get(req.GetId())
	if errors.Is(err, storage.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return toResponse(user), nil
}

// CreateUser создаёт пользователя. Имя и email обязательны, email проверяется.
func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	email := strings.TrimSpace(req.GetEmail())
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return nil, status.Error(codes.InvalidArgument, "email is invalid")
	}

	user, err := s.store.Create(name, email)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return toResponse(user), nil
}

// toResponse превращает доменную модель в ответ gRPC.
// Поле Found оставлено для совместимости с .proto и теперь всегда true:
// отсутствие пользователя сообщается статусом NotFound.
func toResponse(u storage.User) *pb.UserResponse {
	return &pb.UserResponse{
		Id:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Found: true,
	}
}
