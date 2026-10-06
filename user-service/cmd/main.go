package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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
	pb.RegisterUserServiceServer(grpcServer, server.NewUserServer(store))

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Корректная остановка по Ctrl+C или SIGTERM (docker stop):
	// сервер дожидается завершения текущих запросов.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println("user-service останавливается...")
		grpcServer.GracefulStop()
	}()

	log.Println("user-service запущен на :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Println("user-service остановлен")
}
