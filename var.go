package porgs

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"zombiezen.com/go/sqlite/sqlitex"
)

// Context is the root context for the system.
var Context context.Context

// BootConfig holds configuration needed for the system to boot up.
var BootConfig AppBootConfig

// DbConnPool is a pool of SQLite database connections.
var DbConnPool *sqlitex.Pool

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

var ErrNotImplemented = errors.New("not implemented")
var ErrNotFound = errors.New("not found")

const AnonUser = "anon"
const SessionCookieName = "session"
