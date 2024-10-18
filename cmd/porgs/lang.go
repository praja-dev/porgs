package main

import (
	"fmt"
	"github.com/praja-dev/porgs"
	"log/slog"
	"net/http"
)

func handleLang() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		langID := r.PathValue("id")
		slog.Info("core.handleLang", "lang", langID)

		// # Check if the language is supported
		if !porgs.IsLangSupported(langID) {
			porgs.ShowErrorPage(w, r, porgs.ErrorPage{
				Msg:     fmt.Sprintf("Language not supported: %q", langID),
				BackURL: "/",
				Title:   "Unsupported Language",
			})
			return
		}

		// # Save language selection in an HttpOnly cookie
		cookie := http.Cookie{
			Name:     porgs.CookieNameLang,
			Path:     "/",
			Value:    langID,
			MaxAge:   0,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		// # Redirect to the same page that the request came from
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	})
}
