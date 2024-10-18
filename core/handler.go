package core

import (
	"github.com/praja-dev/porgs"
	"net/http"
)

func (p *Plugin) GetHandler() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /{$}", handleRoot())
	mux.Handle("GET /orgs", handleOrgs())
	mux.Handle("GET /org/{id}", handleOrg())

	return mux
}

func handleRoot() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		porgs.RenderView(w, r, porgs.View{Name: "core-root", Title: "Core"})
	})
}
