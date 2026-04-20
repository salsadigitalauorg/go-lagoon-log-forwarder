package logger

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
)

type unescapeHtmlJSONHandler struct {
	writer io.Writer
	opts   *slog.HandlerOptions
}

func (h *unescapeHtmlJSONHandler) Enabled(ctx context.Context, level slog.Level) bool {

	return level >= h.opts.Level.Level()

}

func (h *unescapeHtmlJSONHandler) Handle(ctx context.Context, r slog.Record) error {

	attrs := make(map[string]interface{})

	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})

	encoder := json.NewEncoder(h.writer)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(attrs)
}

func (h *unescapeHtmlJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// dummy
	return h
}

func (h *unescapeHtmlJSONHandler) WithGroup(name string) slog.Handler {
	return h
}
