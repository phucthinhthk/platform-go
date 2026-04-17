//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	"my-project/internal/infrastructure/queryservice"
	"my-project/internal/infrastructure/repository"
	"my-project/internal/interfaces/controller"
	"my-project/internal/usecases"
)

func InitializeHandler(db *gorm.DB) *controller.UserController {
	wire.Build(
		repository.NewMySQLUserRepository,
		queryservice.NewUserQueryService,
		usecases.NewUserUsecase,
		controller.NewUserController,
	)
	return &controller.UserController{}
}
