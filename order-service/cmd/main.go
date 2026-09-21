package main

import (
	"context"
	"log"

	"grpc-demo/order-service/internal/client"
	"grpc-demo/order-service/internal/service"
)

func main() {
	userClient, err := client.NewUserClient("localhost:50051")
	if err != nil {
		log.Fatalf("не удалось подключиться к user-service: %v", err)
	}
	defer userClient.Close()

	orderSvc := service.NewOrderService(userClient)

	order, err := orderSvc.CreateOrder(context.Background(), "8272c4eb-baa6-4118-9d6c-02bd9872fc74", "Клавиатура")
	if err != nil {
		log.Printf("не удалось создать заказ: %v", err)
		return
	}

	log.Printf("заказ создан: %+v", order)
}
