// Package core powers complex organizational structures and the evolving memberships
// of people within those organizations.
//
// A Person (human, company, or system) can hold Membership in multiple organizations (Org)
// simultaneously. However, each person can maintain at most one active Membership
// per organization at any given time.
//
// Membership track status (active, changed, paused, expired) and include start and end timestamps
// for historical recording. When certain fields on a Membership record changes, the current record
// is archived (Membership.Status is set to changed) and a new Membership record is created to
// represent the new situation.
//
// Optionally, a Person can be associated with a porgs.User to enable login capability to PORGS.
//
// A Person can also be optionally linked to a set of other people who act as proxies to the
// original Person. For a human, this could be a trusted friend or family member; for a company,
// it could be a trusted employee or contractor. Systems, considered equivalent to humans,
// includes bots and other automated systems potentially powered by AI.
package core
