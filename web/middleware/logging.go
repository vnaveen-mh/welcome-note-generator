package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vnaveen-mh/welcome-note-generator/web/constants"
	"github.com/vnaveen-mh/welcome-note-generator/web/utils"
)

// responseWriter wraps gin.ResponseWriter to capture response details
type responseWriter struct {
	gin.ResponseWriter
	body []byte
	size int
}

// LoggerMiddleware creates a custom logging middleware using slog
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate or retrieve request ID
		requestID := c.GetHeader(constants.RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
			c.Header(constants.RequestIDHeader, requestID)
		}
		// store request_id in context
		c.Set(constants.RequestIDHeader, requestID)

		isDatastar := utils.IsDatastarRequest(c)

		logger := slog.Default().With(
			slog.String("request_id", requestID),
			slog.String("client_ip", c.ClientIP()),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("query", c.Request.URL.RawQuery),
			slog.Bool("datastar_request", isDatastar),
		)
		utils.SetLogger(c, logger)

		// affects only logging within this func context
		logger = logger.With(slog.String("handler", "LoggingMiddleware"))

		// Create a response writer wrapper to capture response details
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           []byte{},
		}
		c.Writer = writer

		startTime := time.Now()

		// Log request details
		logger.Info("LoggingMiddleware: request_start", slog.Group("request",
			slog.String("user_agent", c.Request.UserAgent()),
			slog.Any("headers", filterHeaders(c.Request.Header)),
			slog.Int64("content_length", c.Request.ContentLength),
			slog.String("protocol", c.Request.Proto),
		),
			slog.Time("timestamp", startTime),
		)

		// Process request
		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// Determine log level based on status code
		statusCode := c.Writer.Status()
		logLevel := slog.LevelInfo
		if statusCode >= 500 {
			logLevel = slog.LevelError
		} else if statusCode >= 400 {
			logLevel = slog.LevelWarn
		}

		// Prepare resp attributes
		responseAttrs := []any{
			slog.Int("status", statusCode),
			slog.Int("size", writer.size),
			slog.Any("headers", filterHeaders(c.Writer.Header())),
		}

		// Add errors if any
		errAttrs := []any{}
		if len(c.Errors) > 0 {
			errorMsgs := make([]string, len(c.Errors))
			for i, err := range c.Errors {
				errorMsgs[i] = err.Error()
			}
			errAttrs = append(errAttrs, slog.Any("errors", errorMsgs))
		}

		// timing related attributes
		timingAttrs := []any{
			slog.Time("start_time", startTime),
			slog.Time("end_time", endTime),
			slog.Duration("latency", latency),
			slog.Float64("latency_ms", float64(latency.Milliseconds())),
		}

		// Log response
		logger.Log(c.Request.Context(), logLevel, "LoggingMiddleware: request_end",
			slog.Group("response", responseAttrs...),
			slog.Group("error_messages", errAttrs...),
			slog.Group("server_timing", timingAttrs...),
		)
	}
}

func (w *responseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

func (w *responseWriter) WriteString(s string) (int, error) {
	size, err := w.ResponseWriter.WriteString(s)
	w.size += size
	return size, err
}

// filterHeaders removes sensitive headers from logging
func filterHeaders(headers map[string][]string) map[string][]string {
	filtered := make(map[string][]string)
	sensitiveHeaders := map[string]bool{
		"Authorization":       true,
		"Cookie":              true,
		"Set-Cookie":          true,
		"X-Api-Key":           true,
		"X-Auth-Token":        true,
		"X-Csrf-Token":        true,
		"Proxy-Authorization": true,
	}

	for key, values := range headers {
		if sensitiveHeaders[key] {
			filtered[key] = []string{"[REDACTED]"}
		} else {
			filtered[key] = values
		}
	}

	return filtered
}
