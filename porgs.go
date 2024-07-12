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

// Plugin interface for PORGS plugin integration.
type Plugin interface {
	// GetName gives the canonical name of this plugin.
	GetName() string

	// GetDependencies declares plugins that this one depends on.
	GetDependencies() []string

	// GetCapabilities lists the capabilities provided by this plugin.
	GetCapabilities() []Capability

	// GetSuggestedRoles recommends role groupings for the plugin's capabilities.
	GetSuggestedRoles() []Role

	// GetFS gives the file system with HTML templates (views/) and static assets (assets/).
	GetFS() embed.FS

	// GetHandler returns the plugin's HTTP router.
	GetHandler() *http.ServeMux
}
