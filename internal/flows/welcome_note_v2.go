package flows

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/vnaveen-mh/welcome-note-generator/internal/types"
)

func RegisterWelcomeNoteFlowV2(g *genkit.Genkit, name string) {
	f := genkit.DefineFlow(g, name, func(ctx context.Context, input *types.WelcomeNoteInput) (string, error) {
		// Implement AI logic here

		// Validate and set defaults
		noteLength := normalizeLength(input.Length)
		lang := normalizeLanguage(input.Language)
		tone := normalizeTone(input.Tone)

		// Build the prompt with tone guidance
		prompt := buildPromptWithTone(input.Occasion, lang, noteLength, tone)
		systemPrompt := buildSystemPromptWithTone()

		resp, err := genkit.Generate(ctx, g,
			ai.WithPrompt(prompt),
			ai.WithSystem(systemPrompt),
		)
		if err != nil {
			return "", err
		}
		return resp.Text(), nil
	})

	SetFlow(name, f)
}

// Build prompt with tone-specific guidance
func buildPromptWithTone(occasion, language, length, tone string) string {

	return fmt.Sprintf(`Create a welcome note based on the details below. Follow the tone, language, and length exactly.
		Occasion: %s
		Language: %s
		Length: %s
		Tone: %s
		`,
		occasion,
		language,
		length,
		tone,
	)
}

// Build system prompt with tone emphasis
func buildSystemPromptWithTone() string {
	prompt := `You are an assistant that writes personalized welcome notes using structured inputs.

Guidelines:
- Use the provided "occasion" as the core theme of the welcome note.
- Adjust your writing style based on the "tone": warm, formal, casual, humorous, professional, or poetic.
- Generate the note in the specified "language".
- Match the requested "length": 
  - short (2–5 sentences),
  - medium (5–10 sentences),
  - long (10+ sentences).
- Do not invent specific factual details that are not implied by the occasion.
- Always stay positive, welcoming, and appropriate for all audiences.
- If the occasion is unclear or unusual, still create a reasonable, friendly welcome note.
- Respond with plain text only. No markdown or bullet points.
`
	/*
	   basePrompt := `You are an expert at writing welcome notes in various languages and styles. You can also generate personalized notes if enough information is available	`

	   toneInstructions := ValidTones

	   	lengthGuide := map[string]string{
	   		"short":  "2-5 sentences",
	   		"medium": "5-10 sentences",
	   		"long":   "more than 10 sentences",
	   	}

	   var sb strings.Builder

	   fmt.Fprintf(&sb, "%s\n\n", basePrompt)
	   fmt.Fprintf(&sb, "Tone Guidelines:\n")

	   	for tone, desc := range toneInstructions {
	   		fmt.Fprintf(&sb, "\t%s: %s\n", tone, desc)
	   	}

	   fmt.Fprintf(&sb, "\nLength Guidelines:\n")

	   	for len, desc := range lengthGuide {
	   		fmt.Fprintf(&sb, "\t%s: %s\n", len, desc)
	   	}

	   personalizationGuidelines := `

	   	Personalize the notes based on information such as guest or host details when available.

	   `
	   fmt.Fprintf(&sb, "\nPersonalization Guidelines:%s\n", personalizationGuidelines)

	   //- Sound positive, welcoming, and suitable for all ages.
	   // - If the provided input contains any negative, hatred, insulting, sarcastic or toxic content, strictly provide a generic a response such that you cannot do it at this time or try again later.

	   otherGuidelines := `
	     - Keep the tone consistent with the value of "tone" mentioned.
	     - Do not invent specific facts you don’t know.
	     - Write as if "host" is welcoming "guest" personally (if values are provided).
	     - Try to infer if the provided input text contains any "occasion" or "event" even if it is a negative one. If it does not refer to any occasion or event, provide a generic message such as "occasion is  or does not make sense".
	     - Respond with plain text only: no markdown, no bullets.
	       `

	   fmt.Fprintf(&sb, "\nOther Guidelines:\n%s\n", otherGuidelines)
	   return sb.String()
	*/
	return prompt
}
