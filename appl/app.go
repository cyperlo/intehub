package appl

import (
	"errors"
	"intehub/pkg/run"
)

var (
	ErrUnsupportedLanguage = errors.New("unsupported language")
)

// AppIns 应用实例包装器
type AppIns struct {
	run.Runnable
}

// NewAppIns 根据语言类型创建应用实例
func NewAppIns(language, code string) (*AppIns, error) {
	var runnable run.Runnable
	var err error

	switch language {
	case "go":
		runnable, err = NewGoAppIns(code)
	case "javascript":
		runnable, err = NewJavaScriptAppIns(code)
	default:
		return nil, ErrUnsupportedLanguage
	}

	if err != nil {
		return nil, err
	}

	return &AppIns{Runnable: runnable}, nil
}
