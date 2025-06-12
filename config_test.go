package logger

import (
	"bytes"
	"log/slog"
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()

	// Test all default values
	tests := []struct {
		name     string
		actual   interface{}
		expected interface{}
	}{
		{"AddSource", cfg.AddSource, true},
		{"ApplicationName", cfg.ApplicationName, ""},
		{"LogChannel", cfg.LogChannel, "LagoonLogs"},
		{"LogHost", cfg.LogHost, ""},
		{"LogPort", cfg.LogPort, 5140},
		{"LogType", cfg.LogType, ""},
		{"MessageVersion", cfg.MessageVersion, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("NewConfig().%s = %v, want %v", tt.name, tt.actual, tt.expected)
			}
		})
	}
}

func TestConfig(t *testing.T) {
	// Save original values
	originalAddSource := addSource
	originalApplicationName := applicationName
	originalLogChannel := logChannel
	originalLogHost := logHost
	originalLogPort := logPort
	originalLogType := logType
	originalMessageVersion := messageVersion

	// Defer restoration
	defer func() {
		addSource = originalAddSource
		applicationName = originalApplicationName
		logChannel = originalLogChannel
		logHost = originalLogHost
		logPort = originalLogPort
		logType = originalLogType
		messageVersion = originalMessageVersion
	}()

	// Test config function
	testCfg := Config{
		AddSource:       false,
		ApplicationName: "test-app",
		LogChannel:      "TestChannel",
		LogHost:         "test.example.com",
		LogPort:         9999,
		LogType:         "test-type",
		MessageVersion:  2,
	}

	// Capture log output
	var logOutput bytes.Buffer
	handler := slog.NewTextHandler(&logOutput, &slog.HandlerOptions{})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	err := config(testCfg)
	if err != nil {
		t.Fatalf("config() returned unexpected error: %v", err)
	}

	// Verify all values were set correctly
	tests := []struct {
		name     string
		actual   interface{}
		expected interface{}
	}{
		{"addSource", addSource, false},
		{"applicationName", applicationName, "test-app"},
		{"logChannel", logChannel, "TestChannel"},
		{"logHost", logHost, "test.example.com"},
		{"logPort", logPort, 9999},
		{"logType", logType, "test-type"},
		{"messageVersion", messageVersion, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("config() did not set %s correctly: got %v, want %v", tt.name, tt.actual, tt.expected)
			}
		})
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	// Save original values
	originalLogHost := logHost
	originalLogType := logType

	// Defer restoration
	defer func() {
		logHost = originalLogHost
		logType = originalLogType
	}()

	// Set valid values
	logHost = "valid.example.com"
	logType = "valid-type"

	// Capture log output
	var logOutput bytes.Buffer
	handler := slog.NewTextHandler(&logOutput, &slog.HandlerOptions{})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	err := validate()
	if err != nil {
		t.Errorf("validate() returned unexpected error with valid config: %v", err)
	}

	// Check that no warnings were logged for logHost
	if bytes.Contains(logOutput.Bytes(), []byte("log.host is not supplied")) {
		t.Error("validate() should not warn when logHost is provided")
	}
}

func TestValidate_EmptyLogHost(t *testing.T) {
	// Save original values
	originalLogHost := logHost
	originalLogType := logType

	// Defer restoration
	defer func() {
		logHost = originalLogHost
		logType = originalLogType
	}()

	// Set test values
	logHost = ""
	logType = "valid-type"

	// Capture log output
	var logOutput bytes.Buffer
	handler := slog.NewTextHandler(&logOutput, &slog.HandlerOptions{})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	err := validate()
	if err != nil {
		t.Errorf("validate() returned unexpected error when only logHost is empty: %v", err)
	}

	// Check that warning was logged for empty logHost
	if !bytes.Contains(logOutput.Bytes(), []byte("log.host is not supplied")) {
		t.Error("validate() should warn when logHost is empty")
	}
}

func TestValidate_EmptyLogType(t *testing.T) {
	// Save original values
	originalLogHost := logHost
	originalLogType := logType

	// Defer restoration
	defer func() {
		logHost = originalLogHost
		logType = originalLogType
	}()

	// Set test values
	logHost = "valid.example.com"
	logType = ""

	err := validate()
	if err == nil {
		t.Error("validate() should return error when logType is empty")
	}

	expectedError := "logType is required"
	if err.Error() != expectedError {
		t.Errorf("validate() returned wrong error: got %q, want %q", err.Error(), expectedError)
	}
}

func TestConfig_WithError(t *testing.T) {
	// Save original values
	originalLogHost := logHost
	originalLogType := logType

	// Defer restoration
	defer func() {
		logHost = originalLogHost
		logType = originalLogType
	}()

	// Test config function with invalid configuration
	testCfg := Config{
		LogType: "", // This should cause an error
		LogHost: "test.example.com",
	}

	err := config(testCfg)
	if err == nil {
		t.Error("config() should return error when logType is empty")
	}

	expectedError := "logType is required"
	if err.Error() != expectedError {
		t.Errorf("config() returned wrong error: got %q, want %q", err.Error(), expectedError)
	}
}

func TestConfig_Integration(t *testing.T) {
	// Test that NewConfig + config works together
	cfg := NewConfig()
	cfg.LogType = "integration-test"
	cfg.LogHost = "integration.example.com"

	// Save original values
	originalAddSource := addSource
	originalApplicationName := applicationName
	originalLogChannel := logChannel
	originalLogHost := logHost
	originalLogPort := logPort
	originalLogType := logType
	originalMessageVersion := messageVersion

	// Defer restoration
	defer func() {
		addSource = originalAddSource
		applicationName = originalApplicationName
		logChannel = originalLogChannel
		logHost = originalLogHost
		logPort = originalLogPort
		logType = originalLogType
		messageVersion = originalMessageVersion
	}()

	// Capture log output
	var logOutput bytes.Buffer
	handler := slog.NewTextHandler(&logOutput, &slog.HandlerOptions{})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	err := config(cfg)
	if err != nil {
		t.Fatalf("config() returned unexpected error: %v", err)
	}

	// Verify defaults were applied correctly
	if addSource != true {
		t.Errorf("Expected addSource to be true, got %v", addSource)
	}
	if logChannel != "LagoonLogs" {
		t.Errorf("Expected logChannel to be 'LagoonLogs', got %s", logChannel)
	}
	if logPort != 5140 {
		t.Errorf("Expected logPort to be 5140, got %d", logPort)
	}

	// Verify custom values were set
	if logType != "integration-test" {
		t.Errorf("Expected logType to be 'integration-test', got %s", logType)
	}
	if logHost != "integration.example.com" {
		t.Errorf("Expected logHost to be 'integration.example.com', got %s", logHost)
	}
}

// Benchmark tests
func BenchmarkNewConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewConfig()
	}
}

func BenchmarkConfig(b *testing.B) {
	cfg := NewConfig()
	cfg.LogType = "benchmark-test"

	// Capture log output to prevent console spam
	var logOutput bytes.Buffer
	handler := slog.NewTextHandler(&logOutput, &slog.HandlerOptions{})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := config(cfg); err != nil {
			b.Fatalf("config() returned error: %v", err)
		}
	}
}
