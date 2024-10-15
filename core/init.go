package core

import (
	"context"
)

func (p *Plugin) GetInit(_ context.Context) error {
	loadData()
	return nil
}
