package appl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"intehub/pkg/params"

	"github.com/dop251/goja"
)

const AppIDJavaScript = "javascript_app"

var (
	ErrInvalidJavaScript = errors.New("invalid JavaScript code")
)

// JavaScriptAppIns JavaScript 应用实例
type JavaScriptAppIns struct {
	program *goja.Program
}

// NewJavaScriptAppIns 创建新的 JavaScript 应用实例
func NewJavaScriptAppIns(code string) (*JavaScriptAppIns, error) {
	if code == "" {
		return nil, ErrInvalidJavaScript
	}

	// 编译代码
	program, err := goja.Compile("script.js", code, false)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidJavaScript, err)
	}

	return &JavaScriptAppIns{program: program}, nil
}

// Run 实现 Runnable 接口
func (j *JavaScriptAppIns) Run(ctx context.Context, input params.Params) (params.Params, error) {
	// 创建新的 VM 实例
	vm := goja.New()

	// 将输入转换为 JS 对象
	inputObj := vm.ToValue(input)
	vm.Set("input", inputObj)

	// 执行脚本
	if _, err := vm.RunProgram(j.program); err != nil {
		return nil, fmt.Errorf("JavaScript execution failed: %w", err)
	}

	// 获取 Run 函数
	runFn, ok := goja.AssertFunction(vm.Get("Run"))
	if !ok {
		return nil, errors.New("Run function not defined")
	}

	// 调用 Run 函数
	result, err := runFn(nil, inputObj)
	if err != nil {
		return nil, fmt.Errorf("Run function error: %w", err)
	}

	// 转换结果
	if result == nil {
		return params.Params{"status": "success"}, nil
	}

	// 尝试将结果转为 JSON 再转回 params
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return params.Params{"status": "success", "result": result.Export()}, nil
	}

	var resultMap map[string]interface{}
	if err := json.Unmarshal(resultJSON, &resultMap); err != nil {
		return params.Params{"status": "success", "result": result.Export()}, nil
	}

	output := make(params.Params)
	for k, v := range resultMap {
		output[k] = v
	}

	return output, nil
}
