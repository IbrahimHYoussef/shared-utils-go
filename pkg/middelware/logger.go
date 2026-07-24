package middelware

import "log/slog"

var logger = slog.Default().With("pkg", "middelware")

func SetLogger(base *slog.Logger) {
	if base == nil {
		logger = slog.Default().With("pkg", "middelware")
		return
	}
	logger = base.With("pkg", "middelware")
}
