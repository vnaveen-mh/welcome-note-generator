package main

import (
	"context"
	"log"

	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/vnaveen-mh/welcome-note-generator/internal/flows"
)

func main() {
	ctx := context.Background()
	g := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
		genkit.WithDefaultModel("googleai/gemini-2.5-flash"),
	)
	if g == nil {
		log.Fatal("error during genkit.Init")
	}

	// Register flows
	flows.RegisterWelcomeNoteFlowV1(g, "welcomeNoteFlowV1")
	flows.RegisterWelcomeNoteFlowV2(g, "welcomeNoteFlowV2")
	flows.RegisterWelcomeNoteFlowV3(g, "welcomeNoteFlowV3")
	flows.RegisterWelcomeNoteFlowSafe(g, "welcomeNoteFlowSafe")
	flows.RegisterWelcomeNoteFlowSmart(g, "welcomeNoteFlowSmart")

	// block indefinitely
	select {}
}
