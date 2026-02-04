package logger

import (
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

var (
	once   sync.Once
	logger *slog.Logger
)

type Config struct {
	LogDir string
}

func Init(cfg Config) error {
	var err error
	once.Do(func() {
		if cfg.LogDir == "" {
			cfg.LogDir = "log"
		}

		if err = os.MkdirAll(cfg.LogDir, 0755); err != nil {
			return
		}

		logFile, openErr := os.OpenFile(filepath.Join(cfg.LogDir, "app.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if openErr != nil {
			err = openErr
			return
		}

		// Create a multi-writer to write to both stdout and file
		// Note: For high performance, you might want to use async writing or just file
		// but for dev/debug, both is useful.
		// The user asked for "saving in datbase all this layers need to logging them ... in new folder log"
		// so file is the priority.

		handler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})

		logger = slog.New(handler)
		slog.SetDefault(logger)
	})
	return err
}

func Get() *slog.Logger {
	if logger == nil {
		// Fallback if not initialized, though Init should be called.
		return slog.Default()
	}
	return logger
}
