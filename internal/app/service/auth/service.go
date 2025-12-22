package auth

import (
	"context"
	"errors"
	ac "intehub/internal/app/context"
	userModel "intehub/internal/app/models/user"
	"intehub/pkg/jwt"
	"time"
)

type service struct {
	userModel userModel.Model
}

func New(userModel userModel.Model) Service {
	return &service{
		userModel: userModel,
	}
}

func (s *service) Login(ctx context.Context, username, password string) (*LoginResponse, error) {
	appCtx := ac.MustGetAppContext(ctx)

	user, err := s.userModel.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if !user.CheckPassword(password) {
		return nil, errors.New("密码错误")
	}

	// 生成JWT token
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role, appCtx.Config.JWT.Secret, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	resp := &LoginResponse{
		Token:    token,
		UserInfo: user,
	}
	return resp, nil
}
