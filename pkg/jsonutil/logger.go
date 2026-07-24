package jsonutil

import "log/slog"

var logger = slog.Default().With("pkg", "jsonutil")

func SetLogger(base *slog.Logger) {
	if base == nil {
		logger = slog.Default().With("pkg", "jsonutil")
		return
	}
	logger = base.With("pkg", "jsonutil")
}
