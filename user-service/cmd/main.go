package main

import (
	"log"
	"net"

	"grpc-demo/user-service/internal/server"
	"grpc-demo/user-service/internal/storage"

	pb "grpc-demo/pb"

	"google.golang.org/grpc"
)

func main() {
	// На старте — хранилище в памяти. Позже можно заменить на PostgreSQL,
	// просто реализовав тот же интерфейс storage.UserStorage поверх БД.
	store := storage.NewInMemoryStorage()

	grpcServer := grpc.NewServer()
	userServer := server.NewUserServer(store)
	pb.RegisterUserServiceServer(grpcServer, userServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("user-service запущен на :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
