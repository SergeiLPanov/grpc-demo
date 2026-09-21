package server

import (
	"context"
	"testing"

	"grpc-demo/user-service/internal/storage"

	pb "grpc-demo/pb"

	"github.com/stretchr/testify/assert"
)

func TestUserServer_GetUser(t *testing.T) {
	store := storage.NewInMemoryStorage()
	created, err := store.Create("Sergei", "sergei@example.com")
	assert.NoError(t, err)

	srv := NewUserServer(store)

	tests := []struct {
		name      string
		requestID string
		wantFound bool
	}{
		{name: "существующий пользователь", requestID: created.ID, wantFound: true},
		{name: "несуществующий пользователь", requestID: "unknown-id", wantFound: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := srv.GetUser(context.Background(), &pb.GetUserRequest{Id: tt.requestID})
			assert.NoError(t, err)
			assert.Equal(t, tt.wantFound, resp.Found)
			if tt.wantFound {
				assert.Equal(t, created.Name, resp.Name)
			}
		})
	}
}

func TestUserServer_CreateUser(t *testing.T) {
	store := storage.NewInMemoryStorage()
	srv := NewUserServer(store)

	resp, err := srv.CreateUser(context.Background(), &pb.CreateUserRequest{
		Name:  "Anna",
		Email: "anna@example.com",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Id)
	assert.Equal(t, "Anna", resp.Name)
	assert.True(t, resp.Found)
}
