package core

import (
	"embed"
)

type Plugin struct{}

func (p *Plugin) GetName() string {
	return "core"
}

func (p *Plugin) GetDependencies() []string {
	return nil
}

//go:embed assets/*
//go:embed views/*.go.html
//go:embed schema.sql
//go:embed seed.sql
var embeddedFS embed.FS

func (p *Plugin) GetFS() embed.FS {
	return embeddedFS
}
