package config

import (
	"encoding/hex"
	"os"
	"strconv"
	"time"
)

// Config
type Config struct {
	Env       string
	Server    ServerConfig
	CSRF      CSRFConfig
	RateLimit RateLimitConfig
}

// ServerConfig
type ServerConfig struct {
	Port string
}

type CSRFConfig struct {
	Key            []byte
	TrustedOrigins []string // Additional trusted origins beyond localhost
}

type RateLimitConfig struct {
	RequestsPerMinute int           // Number of requests allowed per minute per IP
	BurstSize         int           // Burst size for rate limiter
	CleanupInterval   time.Duration // How often to cleanup inactive limiters
	LimiterTTL        time.Duration // How long to keep inactive limiters in memory
}

// Load loads config information from env
func Load() *Config {
	return &Config{
		Env: getEnv("ENV", "DEV"),
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		CSRF: CSRFConfig{
			Key:            getEnvCsrfKey("CSRF_KEY"),
			TrustedOrigins: getEnvSlice("CSRF_TRUSTED_ORIGINS", ","),
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: getEnvInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 30),
			BurstSize:         getEnvInt("RATE_LIMIT_BURST_SIZE", 5),
			CleanupInterval:   getEnvDuration("RATE_LIMIT_CLEANUP_INTERVAL", 5*time.Minute),
			LimiterTTL:        getEnvDuration("RATE_LIMIT_LIMITER_TTL", 15*time.Minute),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvSlice(key string, separator string) []string {
	value := os.Getenv(key)
	if value == "" {
		return []string{}
	}
	// Split by separator and trim whitespace
	parts := []string{}
	for _, part := range splitAndTrim(value, separator) {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

func splitAndTrim(s string, sep string) []string {
	parts := []string{}
	for _, part := range splitString(s, sep) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

func splitString(s string, sep string) []string {
	if s == "" {
		return []string{}
	}
	result := []string{}
	current := ""
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, current)
			current = ""
			i += len(sep) - 1
		} else {
			current += string(s[i])
		}
	}
	result = append(result, current)
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && isSpace(s[start]) {
		start++
	}
	for end > start && isSpace(s[end-1]) {
		end--
	}
	return s[start:end]
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func getEnvCsrfKey(key string) []byte {
	csrfKeyHex := os.Getenv(key)
	if csrfKeyHex == "" {
		panic("CSRF key env variable is required")
	}

	// Decode from hex string to bytes
	csrfKey, err := hex.DecodeString(csrfKeyHex)
	if err != nil {
		panic("invalid CSRF_KEY: must be valid hex string")
	}

	if len(csrfKey) != 32 {
		panic("CSRF_KEY must be 32 bytes")
	}
	return csrfKey
}
