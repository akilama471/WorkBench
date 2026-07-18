package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLoggerLevels(t *testing.T) {
	var buf bytes.Buffer
	log := New(LevelInfo, &buf)

	log.Debug(CategoryApplication, "debug message")
	if buf.Len() != 0 {
		t.Error("debug message should not be logged at INFO level")
	}

	log.Info(CategoryApplication, "info message")
	if !strings.Contains(buf.String(), "info message") {
		t.Error("info message should be logged")
	}
	if !strings.Contains(buf.String(), "INFO") {
		t.Error("log entry should contain INFO level")
	}
}

func TestLoggerCategories(t *testing.T) {
	var buf bytes.Buffer
	log := New(LevelDebug, &buf)

	log.Info(CategoryService, "service message")
	if !strings.Contains(buf.String(), "service") {
		t.Error("log entry should contain category")
	}

	buf.Reset()
	log.Info(CategoryRuntime, "runtime message")
	if !strings.Contains(buf.String(), "runtime") {
		t.Error("log entry should contain runtime category")
	}
}

func TestLoggerKeyValueArgs(t *testing.T) {
	var buf bytes.Buffer
	log := New(LevelInfo, &buf)

	log.Info(CategoryApplication, "test message", "key1", "value1", "key2", 42)
	out := buf.String()
	if !strings.Contains(out, "key1=value1") {
		t.Error("log entry should contain key1=value1")
	}
	if !strings.Contains(out, "key2=42") {
		t.Error("log entry should contain key2=42")
	}
}

func TestLoggerSetLevel(t *testing.T) {
	var buf bytes.Buffer
	log := New(LevelInfo, &buf)

	log.SetLevel(LevelError)

	log.Info(CategoryApplication, "info message")
	if buf.Len() != 0 {
		t.Error("info message should not be logged at ERROR level")
	}

	log.Error(CategoryApplication, "error message")
	if !strings.Contains(buf.String(), "error message") {
		t.Error("error message should be logged")
	}
}

func TestLevelString(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{LevelDebug, "DEBUG"},
		{LevelInfo, "INFO"},
		{LevelWarn, "WARN"},
		{LevelError, "ERROR"},
		{Level(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		if got := tt.level.String(); got != tt.expected {
			t.Errorf("Level(%d).String() = %q, want %q", tt.level, got, tt.expected)
		}
	}
}

func TestPrefixLogger(t *testing.T) {
	var buf bytes.Buffer
	log := New(LevelDebug, &buf)

	prefixed := log.WithPrefix("test")
	prefixed.Info("prefixed message")

	out := buf.String()
	if !strings.Contains(out, "test") {
		t.Error("prefixed log should contain prefix")
	}
	if !strings.Contains(out, "prefixed message") {
		t.Error("prefixed log should contain message")
	}
}
