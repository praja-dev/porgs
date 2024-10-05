package core

import (
	"context"
	"github.com/praja-dev/porgs"
	"net/http"
)

func (p *Plugin) GetHandler(ctx context.Context) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /{$}", handleRoot(ctx))

	return mux
}

func handleRoot(_ context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		porgs.RenderView(w, r, porgs.View{Name: "core-root", Title: "Core"})
	})
}
