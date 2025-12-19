package main

import (
	"intehub/cmd"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func main() {
	// 设置日志
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level: envLevel(),
		}),
	))

	// 创建服务器
	srv, err := cmd.NewServer()
	if err != nil {
		slog.Error("failed to create server", "error", err)
		panic(err)
	}

	// 启动服务器
	srv.MustRun()
}

func envLevel() slog.Level {
	lvlStr, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		lvlStr = "info"
	}

	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(lvlStr)); err != nil {
		lvl = slog.LevelInfo
	}

	return lvl
}
