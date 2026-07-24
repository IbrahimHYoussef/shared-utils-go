package dbutile

import "log/slog"

var logger = slog.Default().With("pkg", "dbutile")

func SetLogger(base *slog.Logger) {
	if base == nil {
		logger = slog.Default().With("pkg", "dbutile")
		return
	}
	logger = base.With("pkg", "dbutile")
}
