package flows

import (
	"context"
	"fmt"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/vnaveen-mh/welcome-note-generator/internal/types"
)

func RegisterWelcomeNoteFlowSafe(g *genkit.Genkit, name string) {

	f := genkit.DefineFlow(g, name, func(ctx context.Context, input *types.WelcomeNoteInput) (*types.SafeWelcomeNoteOutput, error) {
		// 1) Run base generator (V3)
		base, err := genkit.Run(ctx, "run_base_flow", func() (*types.WelcomeNoteV3Output, error) {
			return generateWelcomeNote3(ctx, g, input)
		})
		if err != nil {
			return nil, err
		}

		// 2) Run moderation on the generated note
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

		// 3) Build safe output (V3 + safety)
		out := &types.SafeWelcomeNoteOutput{
			Note:     finalNote,
			Occasion: base.Occasion,
			Language: base.Language,
			Length:   base.Length,
			Tone:     base.Tone,
			Metadata: base.Metadata, // if you keep the V3 metadata

			Blocked:        moderated != nil && moderated.Blocked,
			ModerationNote: "",
			OriginalNote:   originalNote,
		}

		if moderated != nil && moderated.ModerationNote != "" {
			out.ModerationNote = moderated.ModerationNote
		}

		return out, nil
	})

	SetFlow(name, f)
}

func moderateWelcomeNote(ctx context.Context, g *genkit.Genkit, note string) (*types.ModerationResult, error) {
	if strings.TrimSpace(note) == "" {
		return &types.ModerationResult{
			SanitizedNote:  note,
			Blocked:        false,
			ModerationNote: "no content to moderate",
		}, nil
	}

	systemPrompt := `You are a content safety filter that removes toxicity, hate speech, personal attacks,
sexual content, self-harm encouragement, or personally identifiable information.
When issues are found, either redact them or replace them with neutral language
appropriate for a friendly welcome note.
`

	userPrompt := fmt.Sprintf(`
Review the following welcome note for safety issues and return a JSON object.

Rules:
- "sanitizedNote": a safe version of the note with unsafe content removed or rewritten.
  - If the note is already safe, return it unchanged.
  - If only parts are unsafe, rewrite only those parts.
  - If the note is fully blocked, this should be an empty string.
- "blocked": a boolean.
  - Use true only if the content is extremely unsafe and cannot be rewritten safely.
  - Otherwise false.
- "moderationNote": a brief explanation of what was changed or why blocking occurred.
  - Example: "removed insult", "redacted private info", "no issues found".

Output format:
- Respond with a single JSON object only.
- Use exactly these keys: sanitizedNote (string), blocked (boolean), moderationNote (string).
- Do not include any other fields or text.

Welcome note to review:
%s
`, note)

	result, _, err := genkit.GenerateData[types.ModerationResult](ctx, g,
		ai.WithSystem(systemPrompt),
		ai.WithPrompt(userPrompt),
	)
	if err != nil {
		return nil, fmt.Errorf("moderating welcome note: %w", err)
	}

	/*
		if result.SanitizedNote == "" {
			result.SanitizedNote = note
		}
	*/

	return result, nil
}
