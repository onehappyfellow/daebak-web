package controllers

import (
	"net/http"

	"github.com/onehappyfellow/daebak-web/templates"
	"github.com/onehappyfellow/daebak-web/views"
)

func StaticHandler(files ...string) http.HandlerFunc {
	templateFiles := append([]string{"layout.gohtml"}, files...)
	tpl := views.Must(views.ParseFS(templates.FS, templateFiles...))
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}
