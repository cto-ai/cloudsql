package logger

import (
	ctoai "github.com/cto-ai/sdk-go"
)

// Logger imports sdk client
type Logger struct {
	Sdk *ctoai.Sdk
}

// LogError logs error with sdk.track
func (l *Logger) LogError(errorMessage string, errVar error) {
	var tags []string = []string{"anteater", "cloudsql", "track"}
	_ = l.Sdk.Track( // swallowing error because there is nothing can do about sdk.track errorring out
		tags,
		errorMessage,
		map[string]interface{}{
			"err": errVar.Error(),
		},
	)
}

func (l *Logger) LogInfo(message string) {
	var ux = ctoai.NewUx()
	ux.Print(message)
}
