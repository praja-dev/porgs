package core

import "github.com/praja-dev/porgs"

var orgsList = porgs.Capability{
	Name:        "orgs-list",
	Description: "List organizations",
	DashUrlPath: "orgs",
}

var personCreate = porgs.Capability{
	Name:        "person-create",
	Description: "Create person record",
	DashUrlPath: "person/create",
}

func (p *Plugin) GetCapabilities() []porgs.Capability {
	return []porgs.Capability{
		orgsList,
		personCreate,
	}
}

func (p *Plugin) GetSuggestedRoles() []porgs.Role {
	return nil
}
