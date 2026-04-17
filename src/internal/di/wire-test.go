//go:build wireinject
// +build wireinject

package di

import (
	"testing"

	"github.com/google/wire"
	"gorm.io/gorm"

	"my-project/internal/infrastructure/queryservice"
	"my-project/internal/infrastructure/repository"
	"my-project/internal/interfaces/controller"
	"my-project/internal/usecases"
)

type TestHandlerSet struct {
	Handler *controller.UserController
	DB      *gorm.DB
}

func TestInitializeControllers(t *testing.T, db *gorm.DB) (*TestHandlerSet, error) {
	wire.Build(
		repository.NewMySQLUserRepository,
		queryservice.NewUserQueryService,
		usecases.NewUserUsecase,
		controller.NewUserController,
		wire.Struct(new(TestHandlerSet), "*"),
	)
	return nil, nil
}
