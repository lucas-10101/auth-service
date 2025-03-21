package logger

import (
	"context"
	"log/slog"
)

type DummyLogHandler struct{}

func (logger *DummyLogHandler) Enabled(context context.Context, level slog.Level) bool {
	return true
}

func (logger *DummyLogHandler) Handle(context context.Context, record slog.Record) error {
	return nil
}

func (logger *DummyLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return nil
}

func (logger *DummyLogHandler) WithGroup(name string) slog.Handler {
	return nil
}
