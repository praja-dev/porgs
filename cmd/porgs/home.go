package main

import (
	"github.com/praja-dev/porgs"
	"net/http"
)

// vmHome is the view model for the Home screen
type vmHome struct {
	Plugins map[string]porgs.Plugin
}

func handleHome() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		porgs.RenderView(w, r, porgs.View{
			Name:  "main-home",
			Title: "Dashboard | " + porgs.SiteConfig.Title,
			Data: vmHome{
				Plugins: porgs.Plugins,
			},
		})
	})
}
