package logger

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// Logger wraps zerolog.Logger for custom behavior
var Logger zerolog.Logger

// contextKey is used to store the logger in context
type contextKey string

const LoggerKey contextKey = "request-logger"

type ResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.Status = code
	rw.ResponseWriter.WriteHeader(code)
}

func NewLogger(isProd bool) zerolog.Logger {
	if isProd {
		// JSON logs (production)
		Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		// Pretty ConsoleWriter logs (development)
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		// Optional: customize colors or field formats
		output.FormatLevel = func(i any) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		output.FormatMessage = func(i any) string {
			return fmt.Sprintf("%s", i)
		}
		output.FormatFieldName = func(i any) string {
			return fmt.Sprintf("%s:", i)
		}
		output.FormatFieldValue = func(i any) string {
			return fmt.Sprintf("%v", i)
		}

		Logger = zerolog.New(output).With().Timestamp().Logger()
	}

	return Logger
}

// LogError logs a structured error with context
func LogError(ctx context.Context, err error, msg string) {
	GetLogger(ctx).Error().Err(err).Msg(msg)
}

// GetLogger extracts the per-request logger from context
func GetLogger(ctx context.Context) *zerolog.Logger {
	if ctx != nil {
		if l, ok := ctx.Value(LoggerKey).(*zerolog.Logger); ok {
			return l
		}
	}
	return &Logger
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
