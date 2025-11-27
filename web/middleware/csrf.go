package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	"github.com/starfederation/datastar-go/datastar"
	"github.com/vnaveen-mh/welcome-note-generator/web/config"
)

func CsrfMiddlware(cfg *config.Config, r *gin.Engine) http.Handler {
	isProd := strings.EqualFold(cfg.Env, "PROD") || strings.EqualFold(cfg.Env, "PRODUCTION")

	// Build trusted origins based on configured port
	// Always include localhost and 127.0.0.1
	trustedOrigins := []string{
		fmt.Sprintf("localhost:%s", cfg.Server.Port),
		fmt.Sprintf("127.0.0.1:%s", cfg.Server.Port),
	}

	// Add configured trusted origins (can be IPs, domains, etc.)
	for _, origin := range cfg.CSRF.TrustedOrigins {
		// If origin includes port, use as-is, otherwise append configured port
		if containsPort(origin) {
			trustedOrigins = append(trustedOrigins, origin)
		} else {
			trustedOrigins = append(trustedOrigins, fmt.Sprintf("%s:%s", origin, cfg.Server.Port))
		}
	}

	log.Printf("CSRF middleware configured with trusted origins: %v", trustedOrigins)

	// csrf protection middleware
	middleware := csrf.Protect(
		cfg.CSRF.Key,
		csrf.TrustedOrigins(trustedOrigins),
		csrf.Secure(isProd),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.Path("/"),
		csrf.ErrorHandler(http.HandlerFunc(csrfDebugHandler)),
	)
	// wrap gin engine with csrf protection. gin.Engine implements http.Handler interface
	handler := middleware(r)
	return handler
}

func csrfDebugHandler(w http.ResponseWriter, r *http.Request) {
	failureReason := csrf.FailureReason(r)
	r.ParseForm() // Ensure form is parsed

	contentType := r.Header.Get("Content-Type")
	defaultFieldToken := r.FormValue("gorilla.csrf.Token") // Check the default

	headerToken := r.Header.Get("X-CSRF-Token") // The default header

	slog.Error("CSRF Failure",
		slog.String("reason", failureReason.Error()),
		slog.String("host", r.Host),
		slog.String("origin", r.Header.Get("Origin")),
		slog.String("referrer", r.Header.Get("Referer")),
		slog.String("path", r.RequestURI),
		slog.String("content-type", contentType),
		slog.String("X-CSRF-Token from headers", headerToken),
		slog.String("X-CSRF-Token from form fields", defaultFieldToken),
	)

	// Check if this is a Datastar request
	isDatastar := r.Header.Get("Datastar-Request") == "true"
	errorMessage := fmt.Sprintf("CSRF validation failed: %s", failureReason)

	if isDatastar {
		// Get tab name from the request path
		tabName := getTabNameFromPath(r.URL.Path)

		// Set status code
		w.WriteHeader(http.StatusForbidden)

		// Send SSE signal update
		sse := datastar.NewSSE(w, r)
		signals := map[string]interface{}{
			tabName: map[string]interface{}{
				"error":  errorMessage,
				"result": "",
			}}
		sse.MarshalAndPatchSignals(signals)
	} else {
		// For non-Datastar requests, return JSON error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error": errorMessage,
		})
	}
}

// containsPort checks if a string contains a port number
func containsPort(s string) bool {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == ':' {
			return true
		}
		if s[i] == '/' {
			return false
		}
	}
	return false
}
