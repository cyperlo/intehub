package run

import (
	"context"
	"intehub/pkg/params"
)

type Runnable interface {
	Run(ctx context.Context, params params.Params) (output params.Params, err error)
}
