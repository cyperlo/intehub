package context

import (
	"context"

	// "intehub/internal/app/api/v1/auth"

	"intehub/internal/app/config"
	"intehub/internal/app/models"
	"intehub/internal/utils/const/request"
	"log/slog"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type AppContext struct {
	Config *config.Config
	DB     *gorm.DB
	Model  models.Model
}

var ErrAppContextNotSet = errors.New("App Context not set")

func GetAppContext(ctx context.Context) (*AppContext, error) {
	v := ctx.Value(request.AppContext)
	if v == nil {
		return nil, errors.WithStack(ErrAppContextNotSet)
	}
	return v.(*AppContext), nil
}

func FillAppContext(ctx context.Context, appCtx *AppContext) context.Context {
	newCtx := context.WithValue(ctx, request.AppContext, appCtx)
	return newCtx
}

func MustGetAppContext(ctx context.Context) *AppContext {
	appCtx, err := GetAppContext(ctx)
	if err != nil {
		slog.Error(err.Error())
	}
	return appCtx
}
