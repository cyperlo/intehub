package appl

import (
	"context"
	"errors"
	"fmt"
	"intehub/pkg/params"
	"intehub/runtime/yaegi"
	"log/slog"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

const AppIDGo = "go_app"

var (
	ErrInvalidCode       = errors.New("invalid code")
	ErrSignatureMismatch = errors.New("goapp.Run signature mismatch")
)

type GoAppIns struct {
	code    string
	runFunc runFunc
}

type runFunc func(map[string]any) (map[string]any, error)

// NewGoAppIns 创建新的 Go 应用实例
func NewGoAppIns(code string) (*GoAppIns, error) {
	runFunc, err := code2Runnable(code)
	if err != nil {
		return nil, err
	}

	return &GoAppIns{
		code:    code,
		runFunc: runFunc,
	}, nil
}

// Run 实现 Runnable 接口
func (g *GoAppIns) Run(ctx context.Context, input params.Params) (params.Params, error) {
	// 将 params.Params 转换为 map[string]any
	inputMap := make(map[string]any)
	for k, v := range input {
		inputMap[k] = v
	}

	// 执行应用代码
	outputMap, err := g.runFunc(inputMap)
	if err != nil {
		return nil, fmt.Errorf("failed to run app: %w", err)
	}

	// 将结果转换回 params.Params
	output := make(params.Params)
	for k, v := range outputMap {
		output[k] = v
	}

	return output, nil
}

func code2Runnable(code string) (runFunc, error) {
	i := interp.New(interp.Options{})
	if err := loadLibs(i); err != nil {
		return nil, fmt.Errorf("failed to load libs: %w", err)
	}

	_, err := i.Eval(code)
	if err != nil {
		return nil, fmt.Errorf("failed to eval code: %w", err)
	}

	v, err := i.Eval("goapp.Run")
	if err != nil {
		return nil, fmt.Errorf("failed to get goapp.Run: %w", err)
	}
	slog.Debug("app loaded", "app", v)

	rf, ok := v.Interface().(func(map[string]any) (map[string]any, error))
	if !ok {
		return nil, ErrSignatureMismatch
	}

	return rf, nil
}

func loadLibs(i *interp.Interpreter) error {
	if err := i.Use(stdlib.Symbols); err != nil {
		return fmt.Errorf("failed to use stdlib: %w", err)
	}
	if err := i.Use(yaegi.Symbols); err != nil {
		return fmt.Errorf("failed to use plibs: %w", err)
	}
	return nil
}
