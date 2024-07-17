package main

import (
	"github.com/praja-dev/porgs"
	"log/slog"
	"net/http"
)

func handleLoginGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		porgs.RenderView(w, porgs.View{Name: "main-login", Title: "Login | PORGS"})
	})
}

func handleLoginPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			slog.Error("login-post", "err", err)
			porgs.ShowDefaultErrorPage(w)
			return
		}
		_ = r.PostFormValue("username")
		_ = r.PostFormValue("password")

		// TODO: Hash incoming password and compare with the stored version
		// TODO: On success, create a session token and store it in a cookie
		// TODO: On failure, render main-login with username pre-populated and an error message

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
