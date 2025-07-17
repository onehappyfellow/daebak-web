package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"time"

	"github.com/onehappyfellow/daebak-web/context"
	"github.com/onehappyfellow/daebak-web/models"
)

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(patterns[0])
	tpl = tpl.Funcs(
		template.FuncMap{
			"currentUser": func() (*models.User, error) {
				return nil, fmt.Errorf("currentUser not implemented")
			},
			"formatDate": FormatDateLong(),
		},
	)
	tpl, err := tpl.ParseFS(fs, patterns...)
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

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data any) {
	t.htmlTemplate = t.htmlTemplate.Funcs(
		template.FuncMap{
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"formatDate": FormatDateLong(),
		},
	)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.htmlTemplate.Execute(w, data)
	if err != nil {
		fmt.Printf("error executing template: %v", err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}
}

// FormatDateLong converts a time.Time to "Sat, June 14, 2025" format
// Returns the formatted string and nil error for template.FuncMap compatibility
func FormatDateLong() func(time.Time) (string, error) {
	return func(t time.Time) (string, error) {
		return t.Format("Mon, January 2, 2006"), nil
	}
}
