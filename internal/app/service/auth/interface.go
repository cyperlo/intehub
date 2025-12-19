package auth

import (
	"context"
)

type Service interface {
	Login(ctx context.Context, username, password string) (*LoginResponse, error)
}
