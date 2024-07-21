package task

import (
	"github.com/praja-dev/porgs"
	"net/http"
)

func (p *Plugin) GetHandler() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /{$}", handleRoot())

	return mux
}

func handleRoot() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value("user").(porgs.User)
		porgs.RenderView(w, porgs.View{Name: "task-root", Title: "Our Responsibilities", User: u})
	})
}
