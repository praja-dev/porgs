package main

import (
	"github.com/praja-dev/porgs"
	"net/http"
)

func handleRoot() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		porgs.RenderView(w, r, porgs.View{
			Name:  "main-root",
			Title: "Welcome | " + porgs.SiteConfig.Title})
	})
}
