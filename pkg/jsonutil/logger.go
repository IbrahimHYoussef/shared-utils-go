package jsonutil

import "log/slog"

var logger = slog.Default().With("pkg", "jsonutil")

// SetLogger sets the logger used by jsonutil helpers.
//
// Passing nil resets the package logger to slog.Default() with the jsonutil
// package attribute.
func SetLogger(base *slog.Logger) {
	if base == nil {
		logger = slog.Default().With("pkg", "jsonutil")
		return
	}
	logger = base.With("pkg", "jsonutil")
}
