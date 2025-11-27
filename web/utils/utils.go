package utils

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

const (
	loggerKey = "logger"
)

// GetLogger returns the logger from gin.Context
func GetLogger(c *gin.Context) *slog.Logger {
	logger, _ := c.Get(loggerKey)
	if logger != nil {
		return logger.(*slog.Logger)
	}
	return slog.Default()
}

// SetLogger sets the logger within gin.Context
func SetLogger(c *gin.Context, logger *slog.Logger) {
	c.Set(loggerKey, logger)
}
