package logger

import (
	"context"
	"log/slog"
	"regexp"
	"strings"
)

// normalizingHandler wraps an slog.Handler to normalize problematic dynamic fields
// that can cause index explosion in Elasticsearch
type normalizingHandler struct {
	handler slog.Handler
	// patterns defines field paths that should be normalized into arrays
	patterns []*regexp.Regexp
}

// NewNormalizingHandler creates a handler that normalizes dynamic fields
func NewNormalizingHandler(h slog.Handler, pathPatterns []string) slog.Handler {
	patterns := make([]*regexp.Regexp, 0, len(pathPatterns))
	for _, pattern := range pathPatterns {
		if re, err := regexp.Compile(pattern); err == nil {
			patterns = append(patterns, re)
		}
	}
	return &normalizingHandler{
		handler:  h,
		patterns: patterns,
	}
}

func (h *normalizingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *normalizingHandler) Handle(ctx context.Context, r slog.Record) error {
	// Create a new record with normalized attributes
	newRecord := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)

	r.Attrs(func(a slog.Attr) bool {
		newRecord.AddAttrs(h.normalizeAttr([]string{}, a))
		return true
	})

	return h.handler.Handle(ctx, newRecord)
}

func (h *normalizingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	normalizedAttrs := make([]slog.Attr, 0, len(attrs))
	for _, attr := range attrs {
		normalizedAttrs = append(normalizedAttrs, h.normalizeAttr([]string{}, attr))
	}
	return &normalizingHandler{
		handler:  h.handler.WithAttrs(normalizedAttrs),
		patterns: h.patterns,
	}
}

func (h *normalizingHandler) WithGroup(name string) slog.Handler {
	return &normalizingHandler{
		handler:  h.handler.WithGroup(name),
		patterns: h.patterns,
	}
}

// normalizeAttr recursively processes attributes and normalizes matching patterns
func (h *normalizingHandler) normalizeAttr(groups []string, a slog.Attr) slog.Attr {
	// Check if this is a Group that needs normalization
	if a.Value.Kind() == slog.KindGroup {
		groupPath := append(groups, a.Key)
		pathStr := strings.Join(groupPath, ".")

		// Check if this path matches any normalization pattern
		for _, pattern := range h.patterns {
			if pattern.MatchString(pathStr) {
				// Extract all nested attributes and convert to array
				return h.normalizeGroup(a.Key, a.Value.Group())
			}
		}

		// Process group attributes recursively
		attrs := a.Value.Group()
		normalizedAttrs := make([]any, 0, len(attrs))
		for _, attr := range attrs {
			normalizedAttrs = append(normalizedAttrs, h.normalizeAttr(groupPath, attr))
		}
		return slog.Group(a.Key, normalizedAttrs...)
	}

	return a
}

// normalizeGroup converts a group with dynamic keys into an array of key-value pairs
func (h *normalizingHandler) normalizeGroup(key string, attrs []slog.Attr) slog.Attr {
	// Convert to array of maps with "key" and "value" fields
	items := make([]map[string]any, 0, len(attrs))

	for _, attr := range attrs {
		item := make(map[string]any)
		item["key"] = attr.Key

		// Handle nested groups recursively
		if attr.Value.Kind() == slog.KindGroup {
			// For nested groups, flatten them into the value
			item["value"] = groupToMap(attr.Value.Group())
		} else {
			item["value"] = attr.Value.Any()
		}

		items = append(items, item)
	}

	// Return as an array under the original group name
	return slog.Any(key, items)
}

// groupToMap converts nested group attributes to a map for easier consumption
func groupToMap(attrs []slog.Attr) map[string]any {
	result := make(map[string]any)
	for _, attr := range attrs {
		if attr.Value.Kind() == slog.KindGroup {
			result[attr.Key] = groupToMap(attr.Value.Group())
		} else {
			result[attr.Key] = attr.Value.Any()
		}
	}
	return result
}
