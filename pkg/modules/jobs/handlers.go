package jobs

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"
	"shanraq.org/pkg/shanraq"
)

// LogHandler logs the job payload; handy for demos/tests.
func LogHandler(action string) Handler {
	return func(_ context.Context, rt *shanraq.Runtime, job Job) error {
		var payload map[string]any
		_ = json.Unmarshal(job.Payload, &payload)
		rt.Logger.Info("job processed",
			zap.String("action", action),
			zap.String("job_id", job.ID.String()),
			zap.Any("payload", payload))
		return nil
	}
}
