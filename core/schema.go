package core

import "github.com/praja-dev/porgs"

// Person represents an individual person (human, company, or system).
type Person struct {
	ID   int64
	Name string
	Type string

	User porgs.User
}

// Org represents an organization of people.
type Org struct {
	ID   int64
	Name string
	Type string
}

// Membership represents a person's relationship with an organization.
type Membership struct {
	Status string
	Start  string
	End    string

	Person Person
	Org    Org

	Roles []string

	Designation string
	Grade       string
}
