package types

type WelcomeNoteInput struct {
	Occasion string `json:"occasion" form:"occasion" binding:"required" jsonschema:"description=the occasion to generate the welcome note for"`
	Language string `json:"language,omitempty" form:"language" jsonschema:"description=the language of choice for welcome note generation"`
	Length   string `json:"length,omitempty" form:"length" jsonschema:"description=whether the welcome note should be short or medium"`
	Tone     string `json:"tone,omitempty" form:"tone" jsonschema:"description=the tone of the welcome note: formal, casual, warm, humorous, professional, or poetic or insulting or sarcastic"`
}

type WelcomeNoteOutput struct {
	Note     string `json:"note"`
	Occasion string `json:"occasion,omitempty"`
	Language string `json:"language,omitempty"`
	Length   string `json:"length,omitempty"`
	Tone     string `json:"tone,omitempty"`
}

// WelcomeNoteV3Output represents the structured JSON output for V3.
type WelcomeNoteV3Output struct {
	Note     string                `json:"note"`     // final welcome note
	Occasion string                `json:"occasion"` // occasion used by model
	Language string                `json:"language"`
	Length   string                `json:"length"`
	Tone     string                `json:"tone"`
	Metadata WelcomeNoteV3Metadata `json:"metadata"` // nested metadata block
}

// WelcomeNoteV3Metadata provides transparency about how the model interpreted and generated the note.
type WelcomeNoteV3Metadata struct {
	InterpretedOccasion string `json:"interpretedOccasion"` // model's interpretation or normalization
	EffectiveLanguage   string `json:"effectiveLanguage"`   // final actual language used
	EffectiveLength     string `json:"effectiveLength"`     // short | medium | long
	EffectiveTone       string `json:"effectiveTone"`       // warm | formal | poetic | etc.
	Sentiment           string `json:"sentiment"`           // positive | neutral | negative
	Safety              string `json:"safety"`              // safe | needs_review
	Comments            string `json:"comments,omitempty"`  // optional additional notes
}

type SafeWelcomeNoteOutput struct {
	Note     string `json:"note"`
	Occasion string `json:"occasion,omitempty"`
	Language string `json:"language,omitempty"`
	Length   string `json:"length,omitempty"`
	Tone     string `json:"tone,omitempty"`

	// v3-style metadata from generation
	Metadata WelcomeNoteV3Metadata `json:"metadata"`

	// safety info
	Blocked        bool   `json:"blocked"`
	ModerationNote string `json:"moderationNote,omitempty"`
	OriginalNote   string `json:"originalNote,omitempty"` // only set if sanitized
}

type ModerationResult struct {
	SanitizedNote  string `json:"sanitizedNote"`  // empty if no changes
	Blocked        bool   `json:"blocked"`        // true = should not show original
	ModerationNote string `json:"moderationNote"` // short explanation / category
}

type SmartWelcomeFlowOutput struct {
	*SafeWelcomeNoteOutput `json:",inline"` // json tag optional; fields are promoted

	RawDescription string            `json:"rawDescription"`        // what user typed
	ParsedInput    *WelcomeNoteInput `json:"parsedInput,omitempty"` // result of interpretPrompt
}
