package task

import (
	"context"
	"log/slog"
)

func (p *Plugin) GetInit(_ context.Context) error {
	slog.Info("task: loaded")
	return nil
}
