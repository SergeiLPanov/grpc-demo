package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// fakeUsers — заглушка user-service для тестов.
type fakeUsers struct {
	exists bool
	err    error
}

func (f fakeUsers) UserExists(context.Context, string) (bool, error) {
	return f.exists, f.err
}

func TestOrderService_CreateOrder(t *testing.T) {
	errUnavailable := errors.New("user-service недоступен")

	tests := []struct {
		name    string
		users   fakeUsers
		item    string
		wantErr error
	}{
		{name: "успешное создание", users: fakeUsers{exists: true}, item: "Клавиатура"},
		{name: "пользователь не найден", users: fakeUsers{exists: false}, item: "Клавиатура", wantErr: ErrUserNotFound},
		{name: "ошибка user-service", users: fakeUsers{err: errUnavailable}, item: "Клавиатура", wantErr: errUnavailable},
		{name: "пустой товар", users: fakeUsers{exists: true}, item: "   ", wantErr: ErrEmptyItem},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewOrderService(tt.users)

			order, err := svc.CreateOrder(context.Background(), "user-1", tt.item)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, order.ID)
			require.Equal(t, "user-1", order.UserID)
			require.Equal(t, tt.item, order.Item)
		})
	}
}
