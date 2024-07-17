package porgs

import (
	"log/slog"
	"net/http"
)

// View struct holds data to render a view with, and view metadata.
type View struct {
	Name string

	// Title is the display name of the view
	Title string

	// Data is the data to render the view with
	Data interface{}
}

// RenderView renders a view to the response writer.
func RenderView(w http.ResponseWriter, view View) {
	t, ok := Templates[view.Name]
	if !ok {
		slog.Error("render", "view", view, "err", "template not found")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err := t.ExecuteTemplate(w, view.Name, view)
	if err != nil {
		slog.Error("render", "view", view, "err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// ShowDefaultErrorPage shows the default error page with generic data.
func ShowDefaultErrorPage(w http.ResponseWriter) {
	ShowErrorPage(w, ErrorPage{
		Msg:     "There was an unexpected error.",
		Title:   "ShowDefaultErrorPage | PORGS",
		BackURL: "/"})
}

// ShowErrorPage shows an error page with the given error details.
func ShowErrorPage(w http.ResponseWriter, data ErrorPage) {
	RenderView(w, View{Name: "main-error", Title: "ShowDefaultErrorPage | PORGS", Data: data})
}

// ErrorPage struct holds data to render an error page with.
type ErrorPage struct {
	Msg     string
	BackURL string
	Title   string
}
