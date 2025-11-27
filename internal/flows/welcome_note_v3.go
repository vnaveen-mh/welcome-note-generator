package flows

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/vnaveen-mh/welcome-note-generator/internal/types"
)

func RegisterWelcomeNoteFlowV3(g *genkit.Genkit, name string) {
	f := genkit.DefineFlow(g, name, func(ctx context.Context, input *types.WelcomeNoteInput) (*types.WelcomeNoteV3Output, error) {
		// Implement AI logic here
		return generateWelcomeNote3(ctx, g, input)
	})

	SetFlow(name, f)
}

func generateWelcomeNote3(ctx context.Context, g *genkit.Genkit, input *types.WelcomeNoteInput) (*types.WelcomeNoteV3Output, error) {
	input.Length = normalizeLength(input.Length)
	input.Language = normalizeLanguage(input.Language)
	input.Tone = normalizeTone(input.Tone)

	// Build the prompt with tone guidance
	//prompt := buildPromptWithTone(input.Occasion, input.Language, input.Length, input.Tone)
	systemPrompt := buildSystemPromptV3()
	prompt := fmt.Sprintf(
		`Generate the JSON response described in the system prompt using:
Occasion: %s
Language: %s
Length: %s
Tone: %s`,
		input.Occasion, input.Language, input.Length, input.Tone,
	)

	out, _, err := genkit.GenerateData[types.WelcomeNoteV3Output](ctx, g,
		ai.WithPrompt(prompt),
		ai.WithSystem(systemPrompt),
	)
	if err != nil {
		return nil, err
	}

	// Return structured response with metadata
	/*
		return &types.WelcomeNoteOutput{
			Note:     resp.Text(),
			Occasion: input.Occasion,
			Language: input.Language,
			Length:   input.Length,
			Tone:     input.Tone,
		}, nil
	*/
	return out, nil
}

func buildSystemPromptV3() string {
	sysPrompt := `You are an assistant that writes personalized welcome-style notes using structured inputs
and returns both the note and metadata as JSON.

Your role is to generate text, not to enforce content safety policies. A separate moderation
layer will review and, if needed, sanitize your output.

Guidelines:
- Use the provided "occasion" as the main theme of the welcome note.
- Write the note in the specified "language".
- Match the requested "tone" as closely as possible:
  warm, formal, casual, humorous, professional, poetic (or any other tone provided).
- Match the requested "length":
  - short  = about 2–5 sentences
  - medium = about 5–10 sentences
  - long   = about 10+ sentences
- If the occasion is unclear or unusual, still create a reasonable welcome-style message
  and explain your interpretation in the metadata.
- Do not invent specific factual details that are not implied by the input.
- Reflect the sentiment implied by the occasion and tone, even if it is critical,
  frustrated, or darkly humorous. Do not soften or censor strong language that is clearly
  implied by the input just to make it more positive. Safety filtering will be handled
  by another component.

Output format:
- Respond with a single JSON object only, no extra text, no markdown.
- Use this exact structure and key names:

{
  "note": string,                    // the final welcome note
  "occasion": string,                // the occasion you used when writing the note
  "language": string,                // the language you actually used
  "length": string,                  // the length you targeted: short, medium, or long
  "tone": string,                    // the tone you aimed for
  "metadata": {
    "interpretedOccasion": string,   // how you interpreted or normalized the occasion
    "effectiveLanguage": string,     // the final language actually used
    "effectiveLength": string,       // the final length category: short, medium, long
    "effectiveTone": string,         // the final tone you actually wrote in
    "sentiment": "positive" | "neutral" | "negative",
    "safety": "safe" | "needs_review",
    "comments": string               // brief note about any adjustments or concerns
  }
}

- Always produce valid JSON (double quotes around keys and strings, no trailing commas).
`
	return sysPrompt
}
