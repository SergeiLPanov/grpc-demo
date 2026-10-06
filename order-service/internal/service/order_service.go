package service

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("пользователь не найден в user-service")
	ErrEmptyItem    = errors.New("название товара не может быть пустым")
)

// UserChecker — то, что OrderService требует от user-service.
// Интерфейс объявлен здесь, у потребителя, поэтому в тестах
// его легко заменить заглушкой, не поднимая gRPC.
type UserChecker interface {
	UserExists(ctx context.Context, userID string) (bool, error)
}

type Order struct {
	ID     string
	UserID string
	Item   string
}

type OrderService struct {
	mu     sync.RWMutex
	orders map[string]Order
	users  UserChecker
}

func NewOrderService(users UserChecker) *OrderService {
	return &OrderService{
		orders: make(map[string]Order),
		users:  users,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID, item string) (Order, error) {
	if strings.TrimSpace(item) == "" {
		return Order{}, ErrEmptyItem
	}

	exists, err := s.users.UserExists(ctx, userID)
	if err != nil {
		return Order{}, err
	}
	if !exists {
		return Order{}, ErrUserNotFound
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	order := Order{
		ID:     uuid.NewString(),
		UserID: userID,
		Item:   item,
	}
	s.orders[order.ID] = order
	return order, nil
}
