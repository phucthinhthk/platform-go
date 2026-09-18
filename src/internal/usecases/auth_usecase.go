package usecases

import (
	"context"
	"errors"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"platform-go/internal/domain/user"
	"platform-go/internal/infrastructure/security"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthUsecase struct {
	repo user.UserRepository
	jwt  *security.JWTService
}

func NewAuthUsecase(repo user.UserRepository, jwt *security.JWTService) *AuthUsecase {
	return &AuthUsecase{repo: repo, jwt: jwt}
}

func (u *AuthUsecase) Login(ctx context.Context, email, password string) (user.AuthUser, string, string, string, error) {
	account, err := u.repo.FindAuthUserByEmail(ctx, email)
	if err != nil || !account.Active || account.AccountType != "admin" || bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password)) != nil {
		return user.AuthUser{}, "", "", "", ErrInvalidCredentials
	}
	access, err := u.createToken(account, security.AccessTokenType, 24*time.Hour)
	if err != nil {
		return user.AuthUser{}, "", "", "", err
	}
	refresh, err := u.createToken(account, security.RefreshTokenType, 7*24*time.Hour)
	if err != nil {
		return user.AuthUser{}, "", "", "", err
	}
	csrf, err := security.NewCSRFToken()
	if err != nil {
		return user.AuthUser{}, "", "", "", err
	}
	return account, access, refresh, csrf, nil
}

func (u *AuthUsecase) Refresh(ctx context.Context, refreshToken string) (string, string, string, error) {
	claims, err := u.jwt.Parse(refreshToken, security.RefreshTokenType)
	if err != nil {
		return "", "", "", ErrInvalidCredentials
	}
	id, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return "", "", "", ErrInvalidCredentials
	}
	account, err := u.repo.FindAuthUserByID(ctx, id)
	if err != nil || !account.Active || account.AccountType != "admin" {
		return "", "", "", ErrInvalidCredentials
	}
	access, err := u.createToken(account, security.AccessTokenType, 24*time.Hour)
	if err != nil {
		return "", "", "", err
	}
	refresh, err := u.createToken(account, security.RefreshTokenType, 7*24*time.Hour)
	if err != nil {
		return "", "", "", err
	}
	csrf, err := security.NewCSRFToken()
	if err != nil {
		return "", "", "", err
	}
	return access, refresh, csrf, nil
}

func (u *AuthUsecase) CurrentUser(ctx context.Context, subject string) (user.AuthUser, error) {
	id, err := strconv.ParseInt(subject, 10, 64)
	if err != nil {
		return user.AuthUser{}, ErrInvalidCredentials
	}
	account, err := u.repo.FindAuthUserByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) || !account.Active || account.AccountType != "admin" {
		return user.AuthUser{}, ErrInvalidCredentials
	}
	return account, err
}
func (u *AuthUsecase) ParseAccess(token string) (security.Claims, error) {
	return u.jwt.Parse(token, security.AccessTokenType)
}
func (u *AuthUsecase) createToken(a user.AuthUser, kind string, duration time.Duration) (string, error) {
	return u.jwt.Create(security.Claims{Subject: strconv.FormatInt(a.ID, 10), Email: a.Email, AccountType: a.AccountType, TokenType: kind}, duration)
}
