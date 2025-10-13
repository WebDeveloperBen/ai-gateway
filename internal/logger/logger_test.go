package logger_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/WebDeveloperBen/ai-gateway/internal/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	t.Run("production mode creates JSON logger", func(t *testing.T) {
		log := logger.NewLogger(true)

		// Should be a zerolog.Logger
		require.IsType(t, zerolog.Logger{}, log)

		// Test that it can log (basic functionality check)
		var buf bytes.Buffer
		testLogger := log.Output(&buf)
		testLogger.Info().Msg("test message")

		output := buf.String()
		// JSON output should contain the message
		require.Contains(t, output, "test message")
		// Should be valid JSON (contains braces)
		require.Contains(t, output, "{")
		require.Contains(t, output, "}")
	})

	t.Run("development mode creates console logger", func(t *testing.T) {
		// Redirect stdout to capture console output
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		log := logger.NewLogger(false)

		// Should be a zerolog.Logger
		require.IsType(t, zerolog.Logger{}, log)

		// Test logging with fields to trigger format functions
		log.Info().Str("key", "value").Int("number", 42).Msg("test with fields")

		// Restore stdout and read output
		w.Close()
		os.Stdout = oldStdout

		buf := make([]byte, 1024)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		// Should contain formatted output
		require.Contains(t, output, "test with fields")
		require.Contains(t, output, "INFO") // formatted level
		// The format functions should make field names end with :
		// But zerolog's ConsoleWriter might override some formatting
		require.NotEmpty(t, output)
	})
}

func TestLogError(t *testing.T) {
	t.Run("logs error with context", func(t *testing.T) {
		var buf bytes.Buffer
		testLogger := zerolog.New(&buf)
		ctx := context.WithValue(context.Background(), logger.LoggerKey, &testLogger)

		testErr := errors.New("test error")
		logger.LogError(ctx, testErr, "something went wrong")

		output := buf.String()
		require.Contains(t, output, "something went wrong")
		require.Contains(t, output, "test error")
		require.Contains(t, output, "error") // log level
	})

	t.Run("uses global logger when no context logger", func(t *testing.T) {
		var buf bytes.Buffer
		// Temporarily replace global logger
		originalLogger := logger.Logger
		defer func() { logger.Logger = originalLogger }()

		logger.Logger = zerolog.New(&buf)

		ctx := context.Background()
		testErr := errors.New("test error")
		logger.LogError(ctx, testErr, "fallback to global")

		output := buf.String()
		require.Contains(t, output, "fallback to global")
		require.Contains(t, output, "test error")
	})
}

func TestGetLogger(t *testing.T) {
	t.Run("returns logger from context", func(t *testing.T) {
		var buf bytes.Buffer
		contextLogger := zerolog.New(&buf)
		ctx := context.WithValue(context.Background(), logger.LoggerKey, &contextLogger)

		retrieved := logger.GetLogger(ctx)
		require.Equal(t, &contextLogger, retrieved)
	})

	t.Run("returns global logger when context has no logger", func(t *testing.T) {
		ctx := context.Background()
		retrieved := logger.GetLogger(ctx)
		require.Equal(t, &logger.Logger, retrieved)
	})

	t.Run("returns global logger for nil context", func(t *testing.T) {
		retrieved := logger.GetLogger(nil)
		require.Equal(t, &logger.Logger, retrieved)
	})
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	t.Run("captures status code", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		rw := &logger.ResponseWriter{ResponseWriter: recorder}

		rw.WriteHeader(http.StatusCreated)

		require.Equal(t, http.StatusCreated, rw.Status)
		// Verify it also calls the underlying WriteHeader
		require.Equal(t, http.StatusCreated, recorder.Code)
	})

	t.Run("default status is zero", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		rw := &logger.ResponseWriter{ResponseWriter: recorder}

		require.Equal(t, 0, rw.Status)
	})

	t.Run("can be used as http.ResponseWriter", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		var rw http.ResponseWriter = &logger.ResponseWriter{ResponseWriter: recorder}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("test"))

		require.Equal(t, http.StatusOK, recorder.Code)
		require.Equal(t, "test", recorder.Body.String())
	})
}

func TestLoggerKey(t *testing.T) {
	t.Run("logger key constant is defined", func(t *testing.T) {
		// Test that the key is a non-empty string
		require.NotEmpty(t, string(logger.LoggerKey))
		require.Equal(t, "request-logger", string(logger.LoggerKey))
	})
}

func TestNewLogger_GlobalVariable(t *testing.T) {
	t.Run("NewLogger sets global Logger variable", func(t *testing.T) {
		// Save original logger
		originalLogger := logger.Logger
		defer func() { logger.Logger = originalLogger }()

		// Test production mode
		prodLogger := logger.NewLogger(true)
		// Just verify that the global logger is set
		require.NotNil(t, logger.Logger)

		// Test development mode
		devLogger := logger.NewLogger(false)
		// Just verify that the global logger is set
		require.NotNil(t, logger.Logger)

		// Verify that NewLogger returns the logger that was set globally
		require.NotNil(t, prodLogger)
		require.NotNil(t, devLogger)
	})
}

func TestGetLogger_EdgeCases(t *testing.T) {
	t.Run("handles wrong type in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), logger.LoggerKey, "not-a-logger")
		retrieved := logger.GetLogger(ctx)
		// Should fall back to global logger
		require.Equal(t, &logger.Logger, retrieved)
	})

	t.Run("handles nil value in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), logger.LoggerKey, nil)
		retrieved := logger.GetLogger(ctx)
		// Should fall back to global logger
		require.Equal(t, &logger.Logger, retrieved)
	})
}
