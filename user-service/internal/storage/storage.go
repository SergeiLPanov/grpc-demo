package storage

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

// User — доменная модель пользователя.
type User struct {
	ID    string
	Name  string
	Email string
}

// ErrNotFound возвращается, когда пользователь не найден.
var ErrNotFound = errors.New("user not found")

// UserStorage — интерфейс хранилища. Реализовав его поверх PostgreSQL,
// можно заменить InMemoryStorage без изменения остального кода сервиса.
type UserStorage interface {
	Get(id string) (User, error)
	Create(name, email string) (User, error)
}

// InMemoryStorage — простая потокобезопасная реализация в памяти.
// Хороша для старта пет-проекта, потом заменяется на PostgreSQL.
type InMemoryStorage struct {
	mu    sync.RWMutex
	users map[string]User
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		users: make(map[string]User),
	}
}

func (s *InMemoryStorage) Get(id string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return User{}, ErrNotFound
	}
	return user, nil
}

func (s *InMemoryStorage) Create(name, email string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := User{
		ID:    uuid.NewString(),
		Name:  name,
		Email: email,
	}
	s.users[user.ID] = user
	return user, nil
}
