package usecases

import (
	"context"

	"my-project/internal/domain/user"
)

type UserUsecase struct {
	repo user.UserRepository
}

func NewUserUsecase(repo user.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) CreateUser(ctx context.Context, name string) (user.User, error) {
	// Add business logic here if needed
	return u.repo.CreateUser(ctx, name)
}
