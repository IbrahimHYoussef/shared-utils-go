package dbutile

import "log/slog"

var logger = slog.Default().With("pkg", "dbutile")

// SetLogger sets the logger used by dbutile helpers.
//
// Passing nil resets the package logger to slog.Default() with the dbutile
// package attribute.
func SetLogger(base *slog.Logger) {
	if base == nil {
		logger = slog.Default().With("pkg", "dbutile")
		return
	}
	logger = base.With("pkg", "dbutile")
}
