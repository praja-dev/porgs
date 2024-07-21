package main

import (
	"github.com/praja-dev/porgs"
	"net/http"
)

func handleHome() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value("user").(porgs.User)

		porgs.RenderView(w, porgs.View{
			Name:  "main-home",
			Title: "Home | " + porgs.SiteConfig.Title,
			User:  u})
	})
}
