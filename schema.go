package porgs

// User represents a user account in the system.
type User struct {
	Name  string
	Roles []Role
}

// Capability represents a unit of user-facing capability in the system.
// Access control operates at the level of Capabilities.
type Capability struct {
	Name        string
	Description string
	DashUrlPath string
}

// Role is a set of capabilities.
// The set of roles assigned to a User determines what the system permits that user to access.
type Role struct {
	Name         string
	DisplayName  string
	Description  string
	Capabilities []Capability
}

// Access represents an instance where a particular User has accessed a particular Capability.
// This is the basis for the system's audit trail.
type Access struct {
	User       User
	Capability Capability
	Details    string
	Timestamp  string
}

// CxpType represents the type of custom property on an entity—e.g. an organization or person
type CxpType struct {
	Name     string
	Type     string
	Range    string
	Default  string
	Required bool
}
