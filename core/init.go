package core

import (
	"context"
	"log/slog"
)

func (p *Plugin) GetInit(_ context.Context) error {
	loadData()
	slog.Info("core: ready")
	return nil
}
