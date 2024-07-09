package porgs

import "html/template"

// BootConfig holds configuration needed for the system to boot up.
var BootConfig AppBootConfig

// Templates holds all HTML templates for the entire system.
var Templates map[string]*template.Template
