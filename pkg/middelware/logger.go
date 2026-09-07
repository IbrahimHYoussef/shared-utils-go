// Package middelware contains HTTP middleware for request logging,
// authentication, request IDs, and JSON body validation.
package middelware

import "log/slog"

var logger = slog.Default().With("pkg", "middelware")

// SetLogger sets the logger used by middleware helpers.
//
// Passing nil resets the package logger to slog.Default() with the middelware
// package attribute.
func SetLogger(base *slog.Logger) {
	if base == nil {
		logger = slog.Default().With("pkg", "middelware")
		return
	}
	logger = base.With("pkg", "middelware")
}
