package flows

import (
	"strings"

	"github.com/firebase/genkit/go/core/api"
	"github.com/firebase/genkit/go/genkit"
)

// Helper functions for normalization
func normalizeLength(length string) string {
	length = strings.ToLower(strings.TrimSpace(length))
	if length == "short" || length == "medium" || length == "long" {
		return length
	}
	return "short"
}

func normalizeLanguage(language string) string {
	if language == "" {
		return "english"
	}
	return language
}

func normalizeTone(tone string) string {
	tone = strings.ToLower(strings.TrimSpace(tone))
	if _, exists := ValidTones[tone]; exists {
		return tone
	}
	return "warm" // default tone
}

func lookupFlow(g *genkit.Genkit, flowName string) api.Action {
	for _, flow := range genkit.ListFlows(g) {
		if flow.Name() == flowName {
			return flow
		}
	}
	return nil
}
