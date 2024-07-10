package porgs

import "html/template"

// BootConfig holds configuration needed for the system to boot up.
var BootConfig AppBootConfig

// SiteConfig holds the website configuration.
var SiteConfig AppSiteConfig

// Templates holds all HTML templates for the entire system.
var Templates map[string]*template.Template
