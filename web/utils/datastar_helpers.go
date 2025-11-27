package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/starfederation/datastar-go/datastar"
)

// SendSignalUpdateWithError sends an error signal patch to Datastar front end
// If errorMessage is empty, uses the default "error while generating note"
// statusCode is optional and defaults to http.StatusBadRequest
func SendSignalUpdateWithError(c *gin.Context, tabName string, errorMessage string, statusCode ...int) {
	// Use default error message if none provided
	if errorMessage == "" {
		errorMessage = "error while generating note"
	}

	// Use default status code if none provided
	status := http.StatusBadRequest
	if len(statusCode) > 0 {
		status = statusCode[0]
	}

	isDatastar := c.GetHeader("Datastar-Request") == "true"
	if !isDatastar {
		c.JSON(status, gin.H{"error": errorMessage})
		return
	}

	// Set HTTP status code for Datastar requests too
	c.Status(status)
	signals := map[string]interface{}{
		tabName: map[string]interface{}{
			"error":  errorMessage,
			"result": "",
		}}
	SendSignalUpdate(c, signals)
}

// SendSignalUpdate sends a result signal to Datastar
func SendSignalUpdate(c *gin.Context, m map[string]interface{}) {
	sse := datastar.NewSSE(c.Writer, c.Request)
	sse.MarshalAndPatchSignals(m)
}

func IsDatastarRequest(c *gin.Context) bool {
	return c.GetHeader("Datastar-Request") == "true"
}
