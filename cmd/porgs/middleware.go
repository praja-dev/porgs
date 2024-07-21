package main

import (
	"context"
	"log/slog"
	"net/http"
)

// secure middleware check for the id cookie and redirects to /login if not present
func secure(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := r.Cookie("id")
		if err != nil {
			slog.Debug("secure: no cookie")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		u, err := findUserBySessionToken(id.Value)
		if err != nil {
			slog.Error("secure: find user", "err", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), "user", u)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
