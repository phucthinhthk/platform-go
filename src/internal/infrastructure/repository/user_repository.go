package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"platform-go/internal/domain/user"
	"platform-go/internal/infrastructure/dbmodel"
)

type mysqlUserRepository struct {
	db *gorm.DB
}

func (r *mysqlUserRepository) FindAuthUserByEmail(ctx context.Context, email string) (user.AuthUser, error) {
	var u dbmodel.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		return user.AuthUser{}, err
	}
	return authUserFromModel(u), nil
}

func (r *mysqlUserRepository) FindAuthUserByID(ctx context.Context, id int64) (user.AuthUser, error) {
	if id < 1 {
		return user.AuthUser{}, errors.New("invalid user ID")
	}
	var u dbmodel.User
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return user.AuthUser{}, err
	}
	return authUserFromModel(u), nil
}

func authUserFromModel(u dbmodel.User) user.AuthUser {
	var email, passwordHash string
	if u.Email != nil {
		email = *u.Email
	}
	if u.PasswordHash != nil {
		passwordHash = *u.PasswordHash
	}
	return user.AuthUser{ID: int64(u.ID), Name: u.Name, Email: email, PasswordHash: passwordHash, AccountType: u.AccountType, Active: u.Active}
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
