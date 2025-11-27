package logging

import (
	"log/slog"
	"os"
	"sync"
)

var (
	logger *slog.Logger
	once   sync.Once
)

// TBD: should I pass additional options such as log level and attrs
func Init(serviceName, version string) {
	once.Do(func() {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})

		logger = slog.New(handler)
		logger = WithAppInfo(logger, serviceName, version)

		logger.Info("service started")
	})
	slog.SetDefault(logger)
}

func WithAppInfo(logger *slog.Logger, serviceName string, serviceVersion string) *slog.Logger {
	return logger.With(slog.Group("app_info",
		slog.String("app_name", serviceName),
		slog.String("app_version", serviceVersion),
	))
}
