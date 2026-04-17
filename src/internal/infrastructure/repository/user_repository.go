package repository

import (
	"context"

	"gorm.io/gorm"

	"my-project/internal/domain/user"
	"my-project/internal/infrastructure/dbmodel"
)

type mysqlUserRepository struct {
	db *gorm.DB
}

func NewMySQLUserRepository(db *gorm.DB) user.UserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) CreateUser(ctx context.Context, name string) (user.User, error) {
	uData := dbmodel.User{Name: name}
	err := r.db.WithContext(ctx).Create(&uData).Error
	if err != nil {
		return user.User{}, err
	}

	return uData.ToDomain(), nil
}
