package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vnaveen-mh/welcome-note-generator/web/config"
	"github.com/vnaveen-mh/welcome-note-generator/web/utils"
	"golang.org/x/time/rate"
)

// ipLimiterEntry holds the rate limiter and last activity timestamp
type ipLimiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// ipRateLimitStore manages rate limiters per IP with automatic cleanup
type ipRateLimitStore struct {
	mu                sync.RWMutex
	limiters          map[string]*ipLimiterEntry
	requestsPerMinute int
	burstSize         int
	cleanupInterval   time.Duration
	limiterTTL        time.Duration
}

// Global store for IP-based rate limiters
var (
	ipStore     *ipRateLimitStore
	ipStoreOnce sync.Once
)

func newIPRateLimitStore(cfg *config.RateLimitConfig) *ipRateLimitStore {
	s := &ipRateLimitStore{
		limiters:          make(map[string]*ipLimiterEntry),
		requestsPerMinute: cfg.RequestsPerMinute,
		burstSize:         cfg.BurstSize,
		cleanupInterval:   cfg.CleanupInterval,
		limiterTTL:        cfg.LimiterTTL,
	}
	// Start background cleanup goroutine
	go s.cleanupLoop()
	return s
}

// getLimiter retrieves or creates a rate limiter for the given IP
func (s *ipRateLimitStore) getLimiter(ip string) *rate.Limiter {
	// Try read lock first for better performance
	s.mu.RLock()
	entry, exists := s.limiters[ip]
	s.mu.RUnlock()

	if exists {
		// Update last seen time
		s.mu.Lock()
		entry.lastSeen = time.Now()
		s.mu.Unlock()
		return entry.limiter
	}

	// Create new limiter with write lock
	s.mu.Lock()
	defer s.mu.Unlock()

	// Double-check in case another goroutine created it
	entry, exists = s.limiters[ip]
	if exists {
		entry.lastSeen = time.Now()
		return entry.limiter
	}

	// Create new rate limiter
	// Convert requests per minute to requests per second
	limit := rate.Limit(float64(s.requestsPerMinute) / 60.0)
	newLimiter := rate.NewLimiter(limit, s.burstSize)

	s.limiters[ip] = &ipLimiterEntry{
		limiter:  newLimiter,
		lastSeen: time.Now(),
	}

	return newLimiter
}

// cleanupLoop periodically removes inactive limiters
func (s *ipRateLimitStore) cleanupLoop() {
	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanup()
	}
}

// cleanup removes limiters that haven't been used recently
func (s *ipRateLimitStore) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	removed := 0

	for ip, entry := range s.limiters {
		if now.Sub(entry.lastSeen) > s.limiterTTL {
			delete(s.limiters, ip)
			removed++
		}
	}

	if removed > 0 {
		slog.Info("rate limiter cleanup completed",
			slog.Int("removed_ips", removed),
			slog.Int("active_limiters", len(s.limiters)),
		)
	}
}

// getStats returns current rate limiter statistics
func (s *ipRateLimitStore) getStats() (activeIPs int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.limiters)
}

// getTabNameFromPath extracts the tab name from the API path
// e.g., /api/v1/generate -> v1Tab, /api/safe/generate -> safeTab
func getTabNameFromPath(path string) string {
	// Extract the part between /api/ and /generate
	parts := strings.Split(path, "/")
	if len(parts) >= 3 && parts[1] == "api" {
		switch parts[2] {
		case "v1":
			return "v1Tab"
		case "v2":
			return "v2Tab"
		case "v3":
			return "v3Tab"
		case "safe":
			return "safeTab"
		case "smart":
			return "smartTab"
		}
	}
	// Default to empty string if path doesn't match
	return ""
}

// RateLimit is a Gin middleware that enforces per-IP rate limiting
func RateLimit(cfg *config.RateLimitConfig) gin.HandlerFunc {
	// Initialize the global store once
	ipStoreOnce.Do(func() {
		ipStore = newIPRateLimitStore(cfg)
	})

	return func(c *gin.Context) {
		logger := utils.GetLogger(c)

		// Get client IP (handles X-Forwarded-For, X-Real-IP headers)
		clientIP := c.ClientIP()

		logger.Info("rate limit check",
			slog.String("client_ip", clientIP),
		)

		// Get or create rate limiter for this IP
		limiter := ipStore.getLimiter(clientIP)

		// Check if request is allowed
		if !limiter.Allow() {
			logger.Warn("rate limit exceeded",
				slog.String("client_ip", clientIP),
				slog.Int("limit_per_minute", cfg.RequestsPerMinute),
			)

			// Determine tab name from the request path
			tabName := getTabNameFromPath(c.Request.URL.Path)
			errorMessage := fmt.Sprintf("Rate limit exceeded. Maximum %d requests per minute allowed.", cfg.RequestsPerMinute)

			// Use the generic error handler with custom error message and status code
			utils.SendSignalUpdateWithError(c, tabName, errorMessage, http.StatusTooManyRequests)
			c.Abort()
			return
		}

		// Request allowed, continue to next handler
		c.Next()
	}
}
