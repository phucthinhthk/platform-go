package user

import (
	"context"
	"time"
)

type User struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	CreateUser(ctx context.Context, name string) (User, error)
}
