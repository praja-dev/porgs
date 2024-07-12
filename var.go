package porgs

import (
	"html/template"
	"net/http"
)

// BootConfig holds configuration needed for the system to boot up.
var BootConfig AppBootConfig

// SiteConfig holds the website configuration.
var SiteConfig AppSiteConfig

// Plugins holds all the plugins in the system.
var Plugins map[string]Plugin

// Layout holds the parsed HTML template for the layout used by all views.
var Layout *template.Template

// Templates holds a parsed HTML template for each of the views in the system.
var Templates map[string]*template.Template

// Handler is the main HTTP request handler for the system.
var Handler *http.ServeMux
