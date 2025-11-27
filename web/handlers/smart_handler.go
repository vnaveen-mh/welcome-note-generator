package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/firebase/genkit/go/core"
	"github.com/gin-gonic/gin"
	"github.com/vnaveen-mh/welcome-note-generator/internal/flows"
	"github.com/vnaveen-mh/welcome-note-generator/internal/types"
	"github.com/vnaveen-mh/welcome-note-generator/web/utils"
)

type smartInput struct {
	Description string `json:"description" form:"description" binding:"required"`
}

func SmartHandler(c *gin.Context) {
	logger := utils.GetLogger(c)
	logger = logger.With(slog.String("handler", "SmartHandler"))

	logger.Info("http handler begins")
	defer func() {
		logger.Info("http handler ends")
	}()

	isDatastar := utils.IsDatastarRequest(c)

	formInput := smartInput{}
	// This will infer what binder to use depending on the content-type header
	if err := c.ShouldBind(&formInput); err != nil {
		logger.Error("invalid inputs, ShouldBind failed",
			slog.String("error", err.Error()),
			slog.Any("form-input", formInput),
		)

		if !isDatastar {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		utils.SendSignalUpdateWithError(c, "smartTab", "")
		return
	}

	logger.Info("form input", slog.Any("form", formInput))

	val, ok := flows.GetFlow("welcomeNoteFlowSmart")
	if !ok {
		logger.Error("flow does not exist",
			slog.String("error", "flow does not exist in the internal flows store"),
		)
		utils.SendSignalUpdateWithError(c, "smartTab", "")
		return
	}
	flow, ok := val.(*core.Flow[string, *types.SmartWelcomeFlowOutput, struct{}])
	if !ok {
		logger.Error("flow type assertion error",
			slog.String("error", "Flow is not of the right core.Flow type"),
		)
		utils.SendSignalUpdateWithError(c, "smartTab", "")
		return
	}

	// run the flow
	output, err := flow.Run(c.Request.Context(), formInput.Description)

	if err != nil {
		logger.Error("flow.Run returned with error",
			slog.String("error", err.Error()),
		)
		utils.SendSignalUpdateWithError(c, "smartTab", "")
		return
	}

	logger.Info("generated a new note",
		slog.Any("flow.Run output", output),
	)

	resultJson, _ := json.MarshalIndent(output, "", "  ")

	var parsed map[string]interface{}
	if output.ParsedInput != nil {
		parsed = map[string]interface{}{
			"occasion": output.ParsedInput.Occasion,
			"language": output.ParsedInput.Language,
			"length":   output.ParsedInput.Length,
			"tone":     output.ParsedInput.Tone,
		}
	}

	metadata := map[string]interface{}{
		"interpretedOccasion": output.Metadata.InterpretedOccasion,
		"effectiveLanguage":   output.Metadata.EffectiveLanguage,
		"effectiveLength":     output.Metadata.EffectiveLength,
		"effectiveTone":       output.Metadata.EffectiveTone,
		"sentiment":           output.Metadata.Sentiment,
		"safety":              output.Metadata.Safety,
		"comments":            output.Metadata.Comments,
	}

	signals := map[string]interface{}{
		"smartTab": map[string]interface{}{
			"result": map[string]interface{}{
				"note":           output.Note,
				"occasion":       output.Occasion,
				"language":       output.Language,
				"length":         output.Length,
				"tone":           output.Tone,
				"metadata":       metadata,
				"blocked":        output.Blocked,
				"moderationNote": output.ModerationNote,
				"originalNote":   output.OriginalNote,
				"rawDescription": output.RawDescription,
				"parsedInput":    parsed,
			},
			"resultJson": string(resultJson),
			"error":      "",
		},
	}
	if isDatastar {
		utils.SendSignalUpdate(c, signals)
		return
	}
	c.JSON(200, signals["smartTab"])

}
