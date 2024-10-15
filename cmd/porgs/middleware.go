package main

import (
	"context"
	"github.com/praja-dev/porgs"
	"log/slog"
	"net/http"
)

// idUser middleware check for the session cookie and set user in request context.
// If no session cookie is present, a user with name "anon" is set in the context.
func idUser(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := r.Cookie(porgs.SessionCookieName)
		if err != nil {
			ctx := context.WithValue(r.Context(), "user", porgs.User{Name: porgs.AnonUser})
			h.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		u, err := findUserBySessionToken(id.Value)
		if err != nil {
			slog.Error("main.idUser: find user", "err", err)
			u = porgs.User{Name: porgs.AnonUser}
		}

		ctx := context.WithValue(r.Context(), "user", u)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// rejectAnon redirects to /login if the current user is anonymous
func rejectAnon(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, ok := r.Context().Value("user").(porgs.User)
		if !ok {
			u = porgs.User{Name: porgs.AnonUser}
		}

		if u.Name == porgs.AnonUser {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		h.ServeHTTP(w, r)
	})
}
