package server

import (
	"context"
	"errors"

	"grpc-demo/user-service/internal/storage"

	pb "grpc-demo/pb"
)

// UserServer реализует сгенерированный интерфейс pb.UserServiceServer.
// Это и есть та самая "структура, реализующая интерфейс" из плана.
type UserServer struct {
	pb.UnimplementedUserServiceServer // обязательное встраивание для совместимости
	store                             storage.UserStorage
}

func NewUserServer(store storage.UserStorage) *UserServer {
	return &UserServer{store: store}
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := s.store.Get(req.GetId())
	if errors.Is(err, storage.ErrNotFound) {
		return &pb.UserResponse{Found: false}, nil
	}
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Found: true,
	}, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user, err := s.store.Create(req.GetName(), req.GetEmail())
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Found: true,
	}, nil
}
