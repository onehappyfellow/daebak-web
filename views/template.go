package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
)

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl, err := template.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{
		htmlTemplate: tpl,
	}, nil
}

type Template struct {
	htmlTemplate *template.Template
}

func (t Template) Execute(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.htmlTemplate.Execute(w, data)
	if err != nil {
		fmt.Printf("error executing template: %v", err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}
}
