package service

import (
	"errors"
	"intehub/internal/app/model"
	"intehub/pkg/jwt"
	"time"

	"gorm.io/gorm"
)

type AuthService interface {
	Login(username, password string) (string, *model.User, error)
	GetProfile(userID uint) (*model.User, error)
}

type authService struct {
	db     *gorm.DB
	jwtKey string
}

func NewAuthService(db *gorm.DB, jwtKey string) AuthService {
	return &authService{
		db:     db,
		jwtKey: jwtKey,
	}
}

func (s *authService) Login(username, password string) (string, *model.User, error) {
	var user model.User
	err := s.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return "", nil, errors.New("用户名或密码错误")
	}

	if !user.CheckPassword(password) {
		return "", nil, errors.New("用户名或密码错误")
	}

	// 生成 JWT token
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role, s.jwtKey, 24*time.Hour)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}

func (s *authService) GetProfile(userID uint) (*model.User, error) {
	var user model.User
	err := s.db.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
