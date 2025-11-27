package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/firebase/genkit/go/plugins/ollama"
	"github.com/firebase/genkit/go/plugins/server"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	"github.com/vnaveen-mh/welcome-note-generator/internal/flows"
	"github.com/vnaveen-mh/welcome-note-generator/logging"
	"github.com/vnaveen-mh/welcome-note-generator/web/config"
	"github.com/vnaveen-mh/welcome-note-generator/web/handlers"
	"github.com/vnaveen-mh/welcome-note-generator/web/middleware"
	"github.com/vnaveen-mh/welcome-note-generator/web/templates"
)

var (
	appName    = "welcome-note-generator"
	appVersion = "wng-0.1"
)

func RegisterOllamaGptOss(g *genkit.Genkit, ollamaPlugin *ollama.Ollama) {
	model := ollamaPlugin.DefineModel(
		g,
		ollama.ModelDefinition{
			Name: "gpt-oss:latest",
			Type: "generate",
		},
		nil,
	)
	log.Printf("Registered %s model\n", model.Name())
}

func main() {
	ctx := context.Background()

	// Load config
	cfg := config.Load()

	// Initialize the logger and set default logger
	logging.Init(appName, appVersion)

	ollamaPlugin := &ollama.Ollama{
		ServerAddress: "http://localhost:11434",
		Timeout:       120,
	}

	// Initialize Genkit
	g := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{}, ollamaPlugin),
		genkit.WithDefaultModel("googleai/gemini-2.5-flash"),
	)
	if g == nil {
		log.Fatal("error during genkit.Init")
	}

	RegisterOllamaGptOss(g, ollamaPlugin)

	// Register all flows
	flows.RegisterWelcomeNoteFlowV1(g, "welcomeNoteFlowV1")
	flows.RegisterWelcomeNoteFlowV2(g, "welcomeNoteFlowV2")
	flows.RegisterWelcomeNoteFlowV3(g, "welcomeNoteFlowV3")
	flows.RegisterWelcomeNoteFlowSafe(g, "welcomeNoteFlowSafe")
	flows.RegisterWelcomeNoteFlowSmart(g, "welcomeNoteFlowSmart")

	// Set up Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	handler := middleware.CsrfMiddlware(cfg, router)
	//handler = middleware.NormalizeReferer(handler)

	// Helper middleware to add CSRF token to Gin context
	router.Use(func(c *gin.Context) {
		c.Set("csrf_token", csrf.Token(c.Request))
		c.Next()
	})

	router.Use(middleware.Logger())

	// Serve the main page
	router.GET("/", func(c *gin.Context) {
		csrfToken := c.GetString("csrf_token")
		component := templates.Index(csrfToken)
		templ.Handler(component).ServeHTTP(c.Writer, c.Request)
	})

	// API endpoints with IP-based rate limiting
	api := router.Group("/api")
	api.Use(middleware.RateLimit(&cfg.RateLimit))
	{
		api.POST("/v1/generate", handlers.V1Handler)
		api.POST("/v2/generate", handlers.V2Handler)
		api.POST("/v3/generate", handlers.V3Handler)
		api.POST("/safe/generate", handlers.SafeHandler)
		api.POST("/smart/generate", handlers.SmartHandler)
	}

	// Static files (if needed)
	router.Static("/static", "./web/static")

	address := fmt.Sprintf(":%s", cfg.Server.Port)
	slog.Info("🚀 starting web server ", slog.String("address", address))

	mux := http.NewServeMux()
	mux.Handle("/", handler)

	log.Fatal(server.Start(ctx, address, mux))
	//log.Fatal(http.ListenAndServeTLS(address, "cert.pem", "key.pem", handler))
}
