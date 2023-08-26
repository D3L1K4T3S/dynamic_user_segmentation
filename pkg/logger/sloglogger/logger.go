package sloglogger

import (
	"dynamic-user-segmentation/pkg/logger/sloglogger/handle"
	"log/slog"
	"os"
)

func NewLogger(level slog.Level) *slog.Logger {
	options := handle.SlogHandlerOptions{
		Options: &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		},
	}

	return slog.New(options.NewSlogHandler(os.Stdout))
}

func SetLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	default:
		return slog.LevelDebug
	}
}
