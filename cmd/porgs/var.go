package main

import "embed"

//go:embed assets/*
//go:embed layouts/*.go.html
//go:embed views/*.go.html
var embeddedFS embed.FS
