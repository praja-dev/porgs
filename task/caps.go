package task

import "github.com/praja-dev/porgs"

var respList = porgs.Capability{
	Name:        "responsibility-list",
	Description: "List organizational responsibilities",
	DashUrlPath: "",
}

func (p *Plugin) GetCapabilities() []porgs.Capability {
	return []porgs.Capability{
		respList,
	}
}

func (p *Plugin) GetSuggestedRoles() []porgs.Role {
	return nil
}
