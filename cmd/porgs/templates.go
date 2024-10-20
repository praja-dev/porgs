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
		"t": func(lang string, key string) string {
			txt, ok := porgs.SiteConfig.Text[lang]
			if !ok {
				txt = porgs.SiteConfig.Text[porgs.SiteConfig.LangDefault]
			}
			val, ok := txt[key]
			if !ok {
				val = porgs.SiteConfig.Text[porgs.SiteConfig.LangDefault][key]
			}
			if val == "" {
				val = key
			}
			return val
		},
	}

	layout, err := template.New("layout").Funcs(fm).ParseFS(embeddedFS, "layouts/default.go.html")
	if err != nil {
		slog.Error("porgs.getLayoutTemplate", "err", err)
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
		slog.Error("porgs.parseViewTemplates", "err", err)
		os.Exit(1)
	}
	for _, viewFile := range viewFiles {
		viewNameMatches := rgxpViewName.FindStringSubmatch(viewFile)
		if viewNameMatches == nil {
			slog.Error("porgs.parseViewTemplates: get view name",
				"viewFile", viewFile, "err", "incorrect file name")
			os.Exit(1)
		}
		viewName := viewNameMatches[1]

		tp, err := layout.Clone()
		if err != nil {
			slog.Error("porgs.parseViewTemplates: clone layout", "err", err)
			os.Exit(1)
		}
		tp, err = tp.ParseFS(embedFS, viewFile)
		if err != nil {
			slog.Error("porgs.parseViewTemplates: parse", "viewName", viewName, "err", err)
			os.Exit(1)
		}
		tm[viewName] = tp
	}

	return tm
}
