package porgs

// User represents a user account in the system.
type User struct {
	Name  string
	Roles []Role
}

// Capability represents some discrete user-facing capability of the system.
// Access control operates at the level of Capabilities.
type Capability struct {
	Name        string
	Description string
}

// Role is a set of capabilities.
// The set of roles assigned to a User determines what the system permits that user to use.
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
