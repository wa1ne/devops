package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"log/slog"

	"github.com/lmittmann/tint"
)

var logger *slog.Logger

func InitLogger(logDirPath, env string) *slog.Logger {
	if logDirPath == "" {
		logDirPath = "./logs/"
	}

	stdoutHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "15:04:05.000 02/01/2006",
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if err, ok := a.Value.Any().(error); ok {
				aErr := tint.Err(err)
				aErr.Key = a.Key
				return aErr
			}
			return a
		},
	})

	fileLogPath := filepath.Join(logDirPath, "server.stderr")
	file, err := os.OpenFile(fileLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		slog.Error("ошибка при открытии файла логов", "error", err)
		os.Exit(1)
	}
	fileHandler := tint.NewHandler(file, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "15:04:05.000 02/01/2006",
		NoColor:    true,
	})

	var handler slog.Handler
	if env == "prod" {
		handler = fileHandler
	} else if env == "dev" {
		handler = stdoutHandler
	} else {
		slog.Error("неизвестное значение env", "env", env)
		os.Exit(1)
	}

	logger = slog.New(handler)
	return logger
}

func LogError(status int, userErr error, errs ...error) {
	var level slog.Level
	switch {
	case status >= 500:
		level = slog.LevelError
	case status >= 400:
		level = slog.LevelWarn
	default:
		level = slog.LevelInfo
	}

	var attrs []slog.Attr
	for i, err := range errs {
		if err != nil {
			key := fmt.Sprintf("errMsg%d", i+1)
			attrs = append(attrs, slog.Any(key, err))
		}
	}

	logger.LogAttrs(context.Background(), level, userErr.Error(), attrs...)
}
