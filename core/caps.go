package core

import "github.com/praja-dev/porgs"

var respList = porgs.Capability{
	Name:        "person-create",
	Description: "Create person record",
}

func (p *Plugin) GetCapabilities() []porgs.Capability {
	return []porgs.Capability{
		respList,
	}
}

func (p *Plugin) GetSuggestedRoles() []porgs.Role {
	return nil
}
