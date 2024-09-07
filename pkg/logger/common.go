package logger

import (
	"log/slog"
	"os"
)

func InitLog(level slog.Level) {
	slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	})
}
