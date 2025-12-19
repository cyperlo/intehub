package auth

import (
	"context"
	"errors"
	ac "intehub/internal/app/context"
	"intehub/pkg/jwt"
	"log/slog"
	"time"
)

type service struct {
}

func New() Service {
	return &service{}
}

func (s *service) Login(ctx context.Context, username, password string) (*LoginResponse, error) {
	appCtx := ac.MustGetAppContext(ctx)
	user, err := appCtx.Model.UserModel().GetUserByUsername(username)
	if err != nil {
		slog.Error(err.Error())
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
		Token: token,
		User:  user,
	}
	return resp, nil
}
