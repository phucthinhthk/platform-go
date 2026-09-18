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

// AuthUser is the subset of a user record required during authentication.
type AuthUser struct {
	ID           int64
	Name         string
	Email        string
	PasswordHash string
	AccountType  string
	Active       bool
}

type UserRepository interface {
	CreateUser(ctx context.Context, name string) (User, error)
	FindAuthUserByEmail(ctx context.Context, email string) (AuthUser, error)
	FindAuthUserByID(ctx context.Context, id int64) (AuthUser, error)
}
