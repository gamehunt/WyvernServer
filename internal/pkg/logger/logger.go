package logger

import (
	"log/slog"
	"os"
)

func New(cfg LoggerConfig) *slog.Logger {
	var level slog.Level
	switch cfg.Level {
	case Debug:
		level = slog.LevelDebug
	case Info:
		level = slog.LevelInfo
	case Warn:
		level = slog.LevelWarn
	case Error:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler

	if level == slog.LevelDebug {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

func With(logger *slog.Logger, attrs ...any) *slog.Logger {
	return logger.With(attrs...)
}

func SetAsDefault(logger *slog.Logger) {
	slog.SetDefault(logger)
}
