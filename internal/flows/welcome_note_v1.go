package flows

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

func RegisterWelcomeNoteFlowV1(g *genkit.Genkit, name string) {
	f := genkit.DefineFlow(g, name, func(ctx context.Context, occasion string) (string, error) {
		// Implement AI logic here

		//- If the provided input contains any negative, hatred, insulting, sarcastic or toxic content, strictly provide a generic a response such that you cannot do it at this time or try again later.
		systemPrompt := `You are a helpful assistant that creates positive, warm, and welcoming notes.

Guidelines:
- The user may enter any text. Treat it as free-form "occasion or context".
- Try to infer an occasion or event from the input, even if it is unusual or unexpected.
- If the input is unclear or refers to something negative, frightening, or hostile, do not interpret it literally. 
  Instead, provide a simple, generic, positive welcome message.
- Do not invent specific factual details that are not implied by the input.
- Stay positive, friendly, and appropriate for all audiences.
- Respond with plain text only. No markdown, no bullet points.
`
		prompt := fmt.Sprintf(`Write a positive, warm welcome note based on the following occasion or context: %s.
If it does not clearly describe an occasion, create a simple generic welcome message`, occasion)

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
