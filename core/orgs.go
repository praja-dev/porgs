package core

import (
	"context"
	"github.com/praja-dev/porgs"
	"net/http"
)

func handleOrgs(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orgs, err := GetOrgs(ctx)
		if err != nil {
			porgs.ShowDefaultErrorPage(w, r)
			return
		}

		if len(orgs) == 0 {
			porgs.RenderView(w, r, porgs.View{Name: "core-orgs", Title: "Orgs", Data: nil})
			return
		}

		porgs.RenderView(w, r, porgs.View{Name: "core-orgs", Title: "Orgs", Data: orgs})
	})
}

func GetOrgs(_ context.Context) ([]Org, error) {
	var orgs []Org
	return orgs, nil
}
