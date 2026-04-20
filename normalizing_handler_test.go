package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestNormalizingHandler_RenovateBotFields(t *testing.T) {
	var buf bytes.Buffer

	// Create a base JSON handler
	baseHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	// Wrap with normalizing handler
	patterns := []string{
		`^extra\.logData\.hits$`,
		`^extra\..*\.hits$`,
	}
	handler := NewNormalizingHandler(baseHandler, patterns)

	logger := slog.New(handler)

	// Simulate renovatebot logging with dynamic URLs as field names
	logger.Info("Renovate processing",
		slog.Group("extra",
			slog.Group("logData",
				slog.Group("hits",
					slog.Group("https://projects.govcms.gov.au/api/v4/projects/GovCMS%2Fwebhooks/merge_requests/57/notes",
						slog.Int("count", 1),
					),
					slog.Group("https://projects.govcms.gov.au/api/v4/projects/GovCMS%2Fwebhooks/merge_requests/58/notes",
						slog.Int("count", 2),
					),
				),
			),
		),
	)

	// Parse the JSON output
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify the structure
	extra, ok := result["extra"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'extra' field")
	}

	logData, ok := extra["logData"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'extra.logData' field")
	}

	// Check that 'hits' is now an array instead of an object with dynamic keys
	hits, ok := logData["hits"].([]interface{})
	if !ok {
		t.Fatalf("Expected 'hits' to be an array, got: %T", logData["hits"])
	}

	if len(hits) != 2 {
		t.Errorf("Expected 2 items in hits array, got %d", len(hits))
	}

	// Verify structure of array items
	for i, item := range hits {
		hitItem, ok := item.(map[string]interface{})
		if !ok {
			t.Errorf("Item %d: expected object, got %T", i, item)
			continue
		}

		if _, hasKey := hitItem["key"]; !hasKey {
			t.Errorf("Item %d: missing 'key' field", i)
		}

		if _, hasValue := hitItem["value"]; !hasValue {
			t.Errorf("Item %d: missing 'value' field", i)
		}
	}

	t.Logf("Output: %s", buf.String())
}

func TestNormalizingHandler_NormalFieldsUnaffected(t *testing.T) {
	var buf bytes.Buffer

	baseHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	patterns := []string{
		`^extra\.logData\.hits$`,
	}
	handler := NewNormalizingHandler(baseHandler, patterns)

	logger := slog.New(handler)

	// Log normal structured data that shouldn't be normalized
	logger.Info("Normal log",
		slog.Group("extra",
			slog.String("user", "test"),
			slog.Int("count", 42),
			slog.Group("metadata",
				slog.String("version", "1.0"),
			),
		),
	)

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	extra, ok := result["extra"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'extra' field")
	}

	// These should remain as normal fields, not arrays
	if _, ok := extra["user"].(string); !ok {
		t.Error("Expected 'user' to remain a string field")
	}

	if _, ok := extra["count"].(float64); !ok {
		t.Error("Expected 'count' to remain a number field")
	}

	metadata, ok := extra["metadata"].(map[string]interface{})
	if !ok {
		t.Error("Expected 'metadata' to remain an object")
	}

	if _, ok := metadata["version"].(string); !ok {
		t.Error("Expected 'metadata.version' to remain a string field")
	}

	t.Logf("Output: %s", buf.String())
}

func TestNormalizingHandler_EmptyPattern(t *testing.T) {
	var buf bytes.Buffer

	baseHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	// Handler with no patterns should pass through unchanged
	handler := NewNormalizingHandler(baseHandler, []string{})

	logger := slog.New(handler)

	logger.Info("Test", slog.String("key", "value"))

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result["key"] != "value" {
		t.Error("Expected normal key-value structure")
	}
}

func TestNormalizingHandler_WithAttrs(t *testing.T) {
	var buf bytes.Buffer

	baseHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	patterns := []string{
		`^extra\.logData\.hits$`,
	}
	handler := NewNormalizingHandler(baseHandler, patterns)

	logger := slog.New(handler).With(
		slog.Group("extra",
			slog.String("application", "test-app"),
		),
	)

	logger.Info("Test message")

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	extra, ok := result["extra"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'extra' field")
	}

	if extra["application"] != "test-app" {
		t.Error("Expected WithAttrs to work correctly")
	}
}

func TestNormalizingHandler_Enabled(t *testing.T) {
	baseHandler := slog.NewJSONHandler(bytes.NewBuffer(nil), &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})

	handler := NewNormalizingHandler(baseHandler, []string{})

	ctx := context.Background()
	if handler.Enabled(ctx, slog.LevelDebug) {
		t.Error("Debug should be disabled when base handler level is Warn")
	}

	if !handler.Enabled(ctx, slog.LevelWarn) {
		t.Error("Warn should be enabled")
	}

	if !handler.Enabled(ctx, slog.LevelError) {
		t.Error("Error should be enabled")
	}
}
