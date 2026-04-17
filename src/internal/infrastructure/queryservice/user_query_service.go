package queryservice

import (
	"context"

	"gorm.io/gorm"

	"my-project/internal/infrastructure/qsdto"
)

type UserQueryService struct {
	db *gorm.DB
}

func NewUserQueryService(db *gorm.DB) *UserQueryService {
	return &UserQueryService{db: db}
}

func (qs *UserQueryService) FetchUsers(ctx context.Context) ([]qsdto.UserDto, error) {
	var users []qsdto.UserDto
	if err := qs.db.WithContext(ctx).Table("users").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
