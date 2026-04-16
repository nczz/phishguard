package service

import (
	"errors"
	"time"

	"github.com/nczz/phishguard/internal/middleware"
	"github.com/nczz/phishguard/internal/model"
	"github.com/nczz/phishguard/internal/repo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	UserRepo  *repo.UserRepo
	JWTSecret string
}

func (s *AuthService) Login(email, password string) (string, *model.User, error) {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("invalid credentials")
		}
		return "", nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	if !user.IsActive {
		return "", nil, errors.New("account is disabled")
	}
	now := time.Now()
	user.LastLogin = &now
	if err := s.UserRepo.Update(user); err != nil {
		return "", nil, err
	}
	token, err := middleware.GenerateToken(s.JWTSecret, user.ID, user.TenantID, user.Role, user.Email)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (s *AuthService) CreateInitialAdmin(email, password string) error {
	_, err := s.UserRepo.FindByEmail(email)
	if err == nil {
		return errors.New("admin already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	hash, err := s.HashPassword(password)
	if err != nil {
		return err
	}
	return s.UserRepo.Create(&model.User{
		Email:        email,
		Name:         "Platform Admin",
		PasswordHash: hash,
		Role:         "platform_admin",
		IsActive:     true,
		TenantID:     nil,
	})
}
