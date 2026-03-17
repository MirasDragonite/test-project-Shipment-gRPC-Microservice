package logger

import (
	"log/slog"
	"os"
)

var L *slog.Logger

func init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	L = slog.New(handler)
}
