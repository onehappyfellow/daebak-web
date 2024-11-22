package controllers

import (
	"net/http"

	"github.com/onehappyfellow/daebak-web/views"
)

func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}
