package core

import (
	"context"
	"log/slog"
)

func (p *Plugin) GetInit(_ context.Context) error {
	slog.Info("core: loaded")
	return nil
}
