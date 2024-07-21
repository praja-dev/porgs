package main

import (
	"log/slog"
	"net/http"
)

// secure middleware check for the id cookie and redirects to /login if not present
func secure(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("id")
		if err != nil {
			slog.Debug("secure: no cookie")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		h.ServeHTTP(w, r)
	})
}
