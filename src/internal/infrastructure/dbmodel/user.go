package dbmodel

import (
	"time"

	"platform-go/internal/domain/user"
)

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Name         string    `gorm:"size:255;not null"`
	Email        *string   `gorm:"size:255;uniqueIndex"`
	PasswordHash *string   `gorm:"size:255"`
	AccountType  string    `gorm:"size:50;not null;default:admin"`
	Active       bool      `gorm:"not null;default:true"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP;autoUpdateLine"`
}

func (u *User) ToDomain() user.User {
	return user.User{
		ID:        int64(u.ID),
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
