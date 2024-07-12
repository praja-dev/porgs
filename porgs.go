package porgs

import (
	"embed"
	"net/http"
)

// AppBootConfig struct holds configuration required at application boot-up.
type AppBootConfig struct {
	// Host to run the web server on
	Host string

	// Port in the host to run the web server on
	Port int

	// DSN (Data Source Name) for the database connection
	DSN string
}

// AppSiteConfig struct holds the website configuration.
type AppSiteConfig struct {
	Title       string
	Description string
}

// Plugin interface must be implemented by all PORGS plugins
type Plugin interface {
	// GetName is the canonical name of this plugin
	GetName() string

	// GetDependencies returns the list of other plugins required by this plugin
	GetDependencies() []string

	// GetCapabilities returns the list of capabilities provided by this plugin
	GetCapabilities() []Capability

	// GetSuggestedRoles returns the list of roles suggested by this plugin
	GetSuggestedRoles() []Role

	GetFS() embed.FS

	GetHandler() *http.ServeMux
}
