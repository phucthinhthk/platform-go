//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	"platform-go/internal/infrastructure/queryservice"
	"platform-go/internal/infrastructure/repository"
	"platform-go/internal/interfaces/controller"
	"platform-go/internal/usecases"
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
