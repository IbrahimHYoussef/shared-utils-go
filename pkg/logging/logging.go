// Package logging provides slog logger factories for consistent structured
// attributes.
package logging

import "log/slog"

// Factory creates scoped slog loggers from a shared base logger.
type Factory struct {
	base *slog.Logger
}

// NewFactory returns a Factory using base.
//
// If base is nil, slog.Default() is used.
func NewFactory(base *slog.Logger) *Factory {
	if base == nil {
		base = slog.Default()
	}

	return &Factory{base: base}
}

// Package returns a logger with a pkg attribute set to name.
func (f *Factory) Package(name string) *slog.Logger {
	return f.base.With(
		slog.String("pkg", name),
	)
}

// Service returns a logger with a service attribute set to name.
func (f *Factory) Service(name string) *slog.Logger {
	return f.base.With(
		slog.String("service", name),
	)
}

// Repository returns a logger with a repository attribute set to name.
func (f *Factory) Repository(name string) *slog.Logger {
	return f.base.With(
		slog.String("repository", name),
	)
}
