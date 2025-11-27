package handlers

import (
	"log/slog"
	"net/http"

	"github.com/firebase/genkit/go/core"
	"github.com/gin-gonic/gin"
	"github.com/vnaveen-mh/welcome-note-generator/internal/flows"
	"github.com/vnaveen-mh/welcome-note-generator/web/utils"
)

type v1Input struct {
	Occasion string `json:"occasion" form:"occasion" binding:"required"`
}

func V1Handler(c *gin.Context) {
	logger := utils.GetLogger(c)
	logger = logger.With(slog.String("handler", "V1Handler"))

	logger.Info("http handler begins")
	defer func() {
		logger.Info("http handler ends")
	}()
	isDatastar := utils.IsDatastarRequest(c)

	formInput := v1Input{}
	if err := c.ShouldBind(&formInput); err != nil {
		logger.Error("invalid inputs, ShouldBind failed",
			slog.String("error", err.Error()),
			slog.Any("form-input", formInput),
		)

		if !isDatastar {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		utils.SendSignalUpdateWithError(c, "v1Tab", "")
		return
	}

	logger.Info("form input", slog.Any("form", formInput))

	val, ok := flows.GetFlow("welcomeNoteFlowV1")
	if !ok {
		logger.Error("flow does not exist",
			slog.String("error", "flow does not exist in the internal flows store"),
		)
		utils.SendSignalUpdateWithError(c, "v1Tab", "")
		return
	}

	flow, ok := val.(*core.Flow[string, string, struct{}])
	if !ok {
		logger.Error("flow type assertion error",
			slog.String("error", "Flow is not of the right core.Flow type"),
		)
		utils.SendSignalUpdateWithError(c, "v1Tab", "")
		return
	}

	// run the flow
	output, err := flow.Run(c.Request.Context(), formInput.Occasion)

	if err != nil {
		logger.Error("flow.Run returned with error",
			slog.String("error", err.Error()),
		)
		utils.SendSignalUpdateWithError(c, "v1Tab", "")
		return
	}

	logger.Info("generated a new note",
		slog.String("flow.Run output", output),
	)

	signals := map[string]interface{}{
		"v1Tab": map[string]interface{}{
			"result": map[string]interface{}{
				"note": output,
			},
			"error": "",
		},
	}
	if isDatastar {
		utils.SendSignalUpdate(c, signals)
		return
	}
	c.JSON(200, signals["v1Tab"])
}
