package flows

import (
	"context"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/vnaveen-mh/welcome-note-generator/internal/types"
)

func RegisterWelcomeNoteFlowSmart(g *genkit.Genkit, name string) {

	f := genkit.DefineFlow(g, name, func(ctx context.Context, description string) (*types.SmartWelcomeFlowOutput, error) {
		input, err := genkit.Run(ctx, "interpret_description", func() (*types.WelcomeNoteInput, error) {
			return interpretPrompt(ctx, g, description)
		})
		if err != nil {
			return nil, err
		}

		base, err := genkit.Run(ctx, "generate_welcome_note_v3", func() (*types.WelcomeNoteV3Output, error) {
			return generateWelcomeNote3(ctx, g, input)
		})
		if err != nil {
			return nil, err
		}

		moderated, err := genkit.Run(ctx, "moderate_and_sanitize", func() (*types.ModerationResult, error) {
			return moderateWelcomeNote(ctx, g, base.Note)
		})
		if err != nil {
			return nil, err
		}

		finalNote := base.Note
		originalNote := ""

		// Only treat as sanitized if it's non-empty and different from original
		if moderated != nil && moderated.SanitizedNote != "" && moderated.SanitizedNote != base.Note {
			originalNote = base.Note
			finalNote = moderated.SanitizedNote
		}

		safe := &types.SafeWelcomeNoteOutput{
			Note:         finalNote,
			Occasion:     base.Occasion,
			Language:     base.Language,
			Length:       base.Length,
			Tone:         base.Tone,
			Metadata:     base.Metadata,
			Blocked:      moderated != nil && moderated.Blocked,
			OriginalNote: originalNote,
		}
		if moderated != nil && moderated.ModerationNote != "" {
			safe.ModerationNote = moderated.ModerationNote
		}

		// 4) Wrap in Smart output
		out := &types.SmartWelcomeFlowOutput{
			SafeWelcomeNoteOutput: safe,
			RawDescription:        description,
			ParsedInput:           input,
		}
		return out, nil
	})

	SetFlow(name, f)

}

type smartInterpretation struct {
	Occasion string `json:"occasion,omitempty"`
	Language string `json:"language,omitempty"`
	Length   string `json:"length,omitempty"`
	Tone     string `json:"tone,omitempty"`
}

// interpretPrompt uses an LLM to extract a structured WelcomeNoteInput from free-form text.
// Falls back to sensible defaults when the model omits fields.
func interpretPrompt(ctx context.Context, g *genkit.Genkit, description string) (*types.WelcomeNoteInput, error) {
	systemPrompt := `You are an AI that converts free-form descriptions into structured welcome-note inputs.

Your job is NOT to censor or sanitize the user’s text. Your job is to interpret it accurately,
including negative, emotional, sarcastic, or humorous intentions. Safety filtering happens later.

Extract the following fields:
- occasion: What the user is describing (e.g., “welcoming a new hire”, “roasting a bad manager”,
  “sending a sarcastic message”, “celebrating a promotion”).
- language: Infer from the text if clearly indicated; otherwise default to "english".
- tone: Infer from user intent. Valid tones include:
  warm, formal, casual, humorous, professional, poetic,
  AND additional tones when implied: sarcastic, roast, angry, frustrated,
  passive-aggressive, dark-humor, playful, mocking.
- length: Infer short | medium | long.
  Defaults:
    - short = short messages, direct requests, brief sentiments
    - medium = descriptive messages
    - long = highly emotional or detailed requests

Guidelines:
- Do NOT change the user’s meaning.
- Do NOT soften or “nicify” negative sentiments.
- If user clearly requests roasting, criticism, mockery, or negativity, reflect that in "tone".
- If user describes someone negatively (e.g., “my useless boss”), include that in the occasion.
- Never perform safety moderation. That is handled by another component.

Output:
Return only a JSON object:
{
  "occasion": string,
  "language": string,
  "length": string,
  "tone": string
}
`

	userPrompt := "Description: " + description

	result, _, err := genkit.GenerateData[smartInterpretation](ctx, g,
		ai.WithSystem(systemPrompt),
		ai.WithPrompt(userPrompt),
	)
	if err != nil {
		return nil, err
	}

	occ := strings.TrimSpace(result.Occasion)
	if occ == "" {
		occ = strings.TrimSpace(description)
		if occ == "" {
			occ = "welcome event"
		}
	}

	input := &types.WelcomeNoteInput{
		Occasion: occ,
		Language: normalizeLanguage(result.Language),
		Length:   normalizeLength(result.Length),
		Tone:     normalizeTone(result.Tone),
	}

	return input, nil
}
