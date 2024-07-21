package main

import "github.com/praja-dev/porgs"

var authLogin = porgs.Capability{
	Name:        "auth-login",
	Description: "Allow a user to login",
}

var authLogout = porgs.Capability{
	Name:        "auth-logout",
	Description: "Allow an already logged-in user to logout",
}

var authPwdReset = porgs.Capability{
	Name:        "auth-pwd-reset",
	Description: "Allow a user to reset their own password",
}

var authUserCreate = porgs.Capability{
	Name:        "auth-user-create",
	Description: "Create a new user record",
}

var authUserAssignRole = porgs.Capability{
	Name:        "auth-user-assign-role",
	Description: "Assign a role to a user",
}

// getCapabilities gives the list of capabilities exposed by the main package
func getCapabilities() []porgs.Capability {
	return []porgs.Capability{
		authLogin, authLogout, authPwdReset,
		authUserCreate, authUserAssignRole,
	}
}

// getSuggestedRoles gives the roles suggested for organizing the capabilities in main package
func getSuggestedRoles() []porgs.Role {
	return []porgs.Role{
		{

			Name:         "anon",
			DisplayName:  "Anonymous",
			Description:  "As yet unauthenticated user",
			Capabilities: []porgs.Capability{authLogin, authPwdReset},
		},
		{
			Name:         "user",
			DisplayName:  "User",
			Description:  "Already authenticated user",
			Capabilities: []porgs.Capability{authLogout},
		},
		{
			Name:         "admin",
			DisplayName:  "Administrator",
			Description:  "System administrator",
			Capabilities: []porgs.Capability{authUserCreate, authUserAssignRole},
		},
	}
}
