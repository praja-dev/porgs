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
	ID       int64
	Created  int64
	Updated  int64
	ParentID int64
	// SequenceID is the sequence number of this organization within its parent
	SequenceID  int64
	Source      int64
	TypeID      int64
	ExternalID  string
	ExternalSID string
	Name        string
	// Trlx holds translations of standard properties for the organization
	Trlx map[string]OrgProps
	// XProps hold custom properties for the organization
	XProps string
	// XPropsTrlx holds translations of custom properties
	XPropsTrlx string
	SubOrgs    []Org
}

// OrgType represents a type of organization.
type OrgType struct {
	ID          int64
	Created     int64
	Updated     int64
	Name        string
	Description string
	XProps      []porgs.CxpType
}

// OrgProps contains properties common to all organization types
type OrgProps struct{ Name string }

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
