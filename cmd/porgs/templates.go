package main

import (
	"embed"
	"github.com/praja-dev/porgs"
	"html/template"
	"io/fs"
	"log/slog"
	"os"
	"regexp"
)

func getLayoutTemplate() *template.Template {
	fm := template.FuncMap{
		"cfg": func() porgs.AppSiteConfig {
			return porgs.SiteConfig
		},
	}

	layout, err := template.New("layout").Funcs(fm).ParseFS(embeddedFS, "layouts/default.go.html")
	if err != nil {
		slog.Error("templates: parse layouts", "err", err)
		os.Exit(1)
	}

	return layout
}

func getTemplates() map[string]*template.Template {
	tm := parseViewTemplates(embeddedFS, porgs.Layout)

	for _, plugin := range porgs.Plugins {
		pluginTemplates := parseViewTemplates(plugin.GetFS(), porgs.Layout)
		for k, v := range pluginTemplates {
			tm[k] = v
		}
	}

	return tm
}

func parseViewTemplates(embedFS embed.FS, layout *template.Template) map[string]*template.Template {
	tm := make(map[string]*template.Template)

	rgxpViewName := regexp.MustCompile(`views/(.+)\.go\.html`)

	viewFiles, err := fs.Glob(embedFS, "views/*.go.html")
	if err != nil {
		slog.Error("parse views", "err", err)
		os.Exit(1)
	}
	for _, viewFile := range viewFiles {
		viewNameMatches := rgxpViewName.FindStringSubmatch(viewFile)
		if viewNameMatches == nil {
			slog.Error("parse view: incorrect file name", "file", viewFile)
			os.Exit(1)
		}
		viewName := viewNameMatches[1]

		tp, err := layout.Clone()
		if err != nil {
			slog.Error("clone layout", "err", err)
			os.Exit(1)
		}
		tp, err = tp.ParseFS(embedFS, viewFile)
		if err != nil {
			slog.Error("parse view", "view", viewName, "err", err)
			os.Exit(1)
		}
		tm[viewName] = tp
	}

	return tm
}
