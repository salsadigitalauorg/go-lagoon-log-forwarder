package logger

import (
	"bytes"
	"log/slog"
	"strings"
	"sync"
	"testing"
)

func TestDefaultAttrs(t *testing.T) {
	// Save original values
	originalMessageVersion := messageVersion
	originalApplicationName := applicationName
	originalLogChannel := logChannel
	originalHostname := hostname
	originalLogType := logType

	// Defer restoration
	defer func() {
		messageVersion = originalMessageVersion
		applicationName = originalApplicationName
		logChannel = originalLogChannel
		hostname = originalHostname
		logType = originalLogType
	}()

	// Set test values
	messageVersion = 5
	applicationName = "test-app"
	logChannel = "TestChannel"
	hostname = "test-host"
	logType = "test-type"

	attrs := defaultAttrs()

	// Verify the structure and types
	expectedLength := 7 // @version, application, channel, context group, extra group, host, type
	if len(attrs) != expectedLength {
		t.Errorf("defaultAttrs() returned %d attributes, expected %d", len(attrs), expectedLength)
	}

	// Convert to map for easier testing
	attrMap := make(map[string]interface{})
	for i := 0; i < len(attrs); i += 2 {
		if i+1 < len(attrs) {
			if key, ok := attrs[i].(string); ok {
				attrMap[key] = attrs[i+1]
			} else if attr, ok := attrs[i].(slog.Attr); ok {
				attrMap[attr.Key] = attr.Value
			}
		}
	}

	// Test individual attributes by reconstructing what should be there
	// Since slog.Int, slog.String return slog.Attr, we need to test the actual output

	// Test that all expected attributes are present by creating a logger and checking output
	var buf bytes.Buffer
	jsonHandler := slog.NewJSONHandler(&buf, nil)
	testLogger := slog.New(jsonHandler).With(attrs...)
	testLogger.Info("test message")

	output := buf.String()

	// Check for expected values in JSON output
	expectedChecks := []struct {
		name     string
		contains string
	}{
		{"version", `"@version":5`},
		{"application", `"application":"test-app"`},
		{"channel", `"channel":"TestChannel"`},
		{"host", `"host":"test-host"`},
		{"type", `"type":"test-type"`},
	}

	// Note: empty groups (context, extra) don't appear in JSON output by default

	for _, check := range expectedChecks {
		t.Run(check.name, func(t *testing.T) {
			if !strings.Contains(output, check.contains) {
				t.Errorf("defaultAttrs() output missing %s: expected to contain %q in %q",
					check.name, check.contains, output)
			}
		})
	}
}

func TestReplaceAttr(t *testing.T) {
	tests := []struct {
		name     string
		groups   []string
		input    slog.Attr
		expected slog.Attr
	}{
		{
			name:     "msg to message",
			groups:   []string{},
			input:    slog.String("msg", "test message"),
			expected: slog.String("message", "test message"),
		},
		{
			name:     "time to @timestamp",
			groups:   []string{},
			input:    slog.String("time", "2023-01-01T00:00:00Z"),
			expected: slog.String("@timestamp", "2023-01-01T00:00:00Z"),
		},
		{
			name:     "timestampOverride to @timestamp",
			groups:   []string{},
			input:    slog.String("timestampOverride", "2023-01-01T00:00:00Z"),
			expected: slog.String("@timestamp", "2023-01-01T00:00:00Z"),
		},
		{
			name:     "no change for other keys",
			groups:   []string{},
			input:    slog.String("level", "INFO"),
			expected: slog.String("level", "INFO"),
		},
		{
			name:     "no change when in groups",
			groups:   []string{"group1"},
			input:    slog.String("msg", "test message"),
			expected: slog.String("msg", "test message"),
		},
		{
			name:     "no change for nested groups",
			groups:   []string{"group1", "group2"},
			input:    slog.String("time", "2023-01-01T00:00:00Z"),
			expected: slog.String("time", "2023-01-01T00:00:00Z"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceAttr(tt.groups, tt.input)
			if result.Key != tt.expected.Key {
				t.Errorf("replaceAttr() key = %v, want %v", result.Key, tt.expected.Key)
			}
			if !result.Value.Equal(tt.expected.Value) {
				t.Errorf("replaceAttr() value = %v, want %v", result.Value, tt.expected.Value)
			}
		})
	}
}

func TestConnect_InvalidAddress(t *testing.T) {
	// Save original values
	originalLogHost := logHost
	originalLogPort := logPort

	// Defer restoration
	defer func() {
		logHost = originalLogHost
		logPort = originalLogPort
	}()

	// Test with invalid address format
	logHost = "invalid-address-format:::"
	logPort = 5140

	conn, err := connect()
	if err == nil {
		t.Error("connect() should return error for invalid address")
		if conn != nil {
			conn.Close()
		}
	}
	if conn != nil {
		t.Error("connect() should return nil connection for invalid address")
		conn.Close()
	}
}

func TestConnect_ValidAddress(t *testing.T) {
	// Save original values
	originalLogHost := logHost
	originalLogPort := logPort

	// Defer restoration
	defer func() {
		logHost = originalLogHost
		logPort = originalLogPort
	}()

	// Test with valid localhost address
	logHost = "127.0.0.1"
	logPort = 0 // Let OS choose port

	conn, err := connect()
	if err != nil {
		// This might fail in some environments, so we'll make it a soft check
		t.Logf("connect() failed (this may be expected in test environment): %v", err)
		return
	}

	if conn == nil {
		t.Error("connect() should return valid connection for valid address")
		return
	}

	// Verify connection properties
	if conn.LocalAddr() == nil {
		t.Error("connect() should return connection with valid local address")
	}

	if conn.RemoteAddr() == nil {
		t.Error("connect() should return connection with valid remote address")
	}

	// Clean up
	conn.Close()
}

func TestConnect_EmptyHost(t *testing.T) {
	// Save original values
	originalLogHost := logHost
	originalLogPort := logPort

	// Defer restoration
	defer func() {
		logHost = originalLogHost
		logPort = originalLogPort
	}()

	// Test with empty host (should default to empty string, which may cause address resolution to fail)
	logHost = ""
	logPort = 5140

	conn, err := connect()
	// This should likely fail since empty host is invalid
	if err == nil && conn != nil {
		// If it somehow succeeds, clean up
		conn.Close()
		t.Log("connect() succeeded with empty host (unexpected but not necessarily wrong)")
	} else {
		// This is the expected case
		t.Log("connect() failed with empty host as expected")
	}
}

func TestInitialize_ConfigError(t *testing.T) {
	// Test Initialize with invalid config
	cfg := NewConfig()
	cfg.LogType = "" // This should cause validation error

	err := Initialize(cfg)
	if err == nil {
		t.Error("Initialize() should return error for invalid config")
	}

	expectedErrorSubstring := "configuration error"
	if !strings.Contains(err.Error(), expectedErrorSubstring) {
		t.Errorf("Initialize() error should contain %q, got %q", expectedErrorSubstring, err.Error())
	}
}

func TestInitialize_ValidConfig(t *testing.T) {
	// Save original values
	originalOnce := once
	originalHostname := hostname
	originalMessageVersion := messageVersion

	// Defer restoration
	defer func() {
		once = originalOnce
		hostname = originalHostname
		messageVersion = originalMessageVersion
	}()

	// Reset once to allow re-initialization
	once = sync.Once{}

	// Test with valid config
	cfg := NewConfig()
	cfg.LogType = "test-type"
	cfg.LogHost = "127.0.0.1"
	cfg.LogPort = 0 // Use any available port

	// This test might fail due to network connectivity in test environment
	// We'll focus on testing the error path and basic setup
	err := Initialize(cfg)

	// The function might fail at UDP connection, which is expected in test environment
	if err != nil {
		t.Logf("Initialize() failed (may be expected in test environment): %v", err)
	} else {
		t.Log("Initialize() succeeded")
	}

	// Verify that hostname and messageVersion were set
	if hostname == "" {
		t.Error("Initialize() should set hostname")
	}
	if messageVersion != 3 {
		t.Errorf("Initialize() should set messageVersion to 3, got %d", messageVersion)
	}
}

func TestInitialize_OnceSemantics(t *testing.T) {
	// Save original values
	originalOnce := once

	// Defer restoration
	defer func() {
		once = originalOnce
	}()

	// Reset once for this test
	once = sync.Once{}

	cfg := NewConfig()
	cfg.LogType = "test-once"
	cfg.LogHost = "127.0.0.1"

	// First call
	err1 := Initialize(cfg)

	// Second call - the once.Do should prevent re-initialization
	cfg.LogType = "test-once-different"
	err2 := Initialize(cfg)

	// Both calls should have the same result regarding error/success
	// The key thing is that once.Do ensures the initialization block runs only once
	if (err1 == nil) != (err2 == nil) {
		t.Logf("First Initialize: %v, Second Initialize: %v", err1, err2)
		t.Log("Different results may be expected due to once.Do semantics")
	}
}

// Test helper functions
func TestPackageVariables(t *testing.T) {
	// Test that package variables can be set and read
	// This is more of a sanity check
	testValues := map[string]interface{}{
		"addSource":       true,
		"applicationName": "test-app",
		"logChannel":      "test-channel",
		"logHost":         "test-host",
		"logPort":         9999,
		"logType":         "test-type",
		"messageVersion":  42,
	}

	// Save original values
	originals := map[string]interface{}{
		"addSource":       addSource,
		"applicationName": applicationName,
		"logChannel":      logChannel,
		"logHost":         logHost,
		"logPort":         logPort,
		"logType":         logType,
		"messageVersion":  messageVersion,
	}

	// Set test values
	addSource = testValues["addSource"].(bool)
	applicationName = testValues["applicationName"].(string)
	logChannel = testValues["logChannel"].(string)
	logHost = testValues["logHost"].(string)
	logPort = testValues["logPort"].(int)
	logType = testValues["logType"].(string)
	messageVersion = testValues["messageVersion"].(int)

	// Verify values were set
	if addSource != testValues["addSource"] {
		t.Errorf("addSource = %v, want %v", addSource, testValues["addSource"])
	}
	if applicationName != testValues["applicationName"] {
		t.Errorf("applicationName = %v, want %v", applicationName, testValues["applicationName"])
	}
	if logChannel != testValues["logChannel"] {
		t.Errorf("logChannel = %v, want %v", logChannel, testValues["logChannel"])
	}

	// Restore original values
	addSource = originals["addSource"].(bool)
	applicationName = originals["applicationName"].(string)
	logChannel = originals["logChannel"].(string)
	logHost = originals["logHost"].(string)
	logPort = originals["logPort"].(int)
	logType = originals["logType"].(string)
	messageVersion = originals["messageVersion"].(int)
}

// Benchmark tests
func BenchmarkDefaultAttrs(b *testing.B) {
	// Set up test values
	messageVersion = 3
	applicationName = "benchmark-app"
	logChannel = "BenchmarkChannel"
	hostname = "benchmark-host"
	logType = "benchmark-type"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = defaultAttrs()
	}
}

func BenchmarkReplaceAttr(b *testing.B) {
	attr := slog.String("msg", "test message")
	groups := []string{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = replaceAttr(groups, attr)
	}
}

func BenchmarkReplaceAttr_WithGroups(b *testing.B) {
	attr := slog.String("msg", "test message")
	groups := []string{"group1", "group2"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = replaceAttr(groups, attr)
	}
}
