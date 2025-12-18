package yaegi

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type Runtime struct {
	interpreter *interp.Interpreter
}

func NewRuntime() *Runtime {
	return &Runtime{}
}

// Run 执行 Go 代码
func (r *Runtime) Run(ctx context.Context, code string) (string, error) {
	// 创建新的解释器
	i := interp.New(interp.Options{})

	// 导入标准库
	if err := i.Use(stdlib.Symbols); err != nil {
		return "", fmt.Errorf("failed to load stdlib: %w", err)
	}

	// 导入自定义符号
	if err := i.Use(Symbols); err != nil {
		return "", fmt.Errorf("failed to load custom symbols: %w", err)
	}

	// 捕获输出
	var output bytes.Buffer
	oldStdout := os.Stdout
	pipeReader, pipeWriter, _ := os.Pipe()
	os.Stdout = pipeWriter

	// 创建完成通道
	done := make(chan error, 1)

	go func() {
		// 执行代码
		_, err := i.Eval(code)
		pipeWriter.Close()
		done <- err
	}()

	// 等待执行完成或超时
	select {
	case <-ctx.Done():
		os.Stdout = oldStdout
		return "", errors.New("execution timeout")
	case err := <-done:
		// 读取输出
		io.Copy(&output, pipeReader)
		os.Stdout = oldStdout

		if err != nil {
			return output.String(), err
		}
		return output.String(), nil
	}
}
