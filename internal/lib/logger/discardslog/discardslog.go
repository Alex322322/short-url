package discardslog

import (
	"context"
	"log/slog"
)

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardsHandler())
}

type DiscardsHandler struct{}

func NewDiscardsHandler() *DiscardsHandler {
	return &DiscardsHandler{}
}

func (h *DiscardsHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (h *DiscardsHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *DiscardsHandler) WithGroup(_ string) slog.Handler {
	return h
}	

func (h *DiscardsHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}