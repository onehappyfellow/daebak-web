package views

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

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
			"currentUser": func() *models.User {
				return nil
			},
			"toastMessages": func() []Toast {
				return nil
			},
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

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data any, errs ...error) {
	tpl, err := t.htmlTemplate.Clone()
	if err != nil {
		fmt.Printf("error executing template: %v", err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
		return
	}
	tpl = tpl.Funcs(
		template.FuncMap{
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"toastMessages": func() []Toast {
				var toasts []Toast
				for _, err := range errs {
					var pubErr public
					if errors.As(err, &pubErr) {
						toasts = append(toasts, Toast{
							Text: pubErr.Public(),
							Type: "error",
						})
					} else {
						toasts = append(toasts, Toast{
							Text: "Something went wrong.",
							Type: "error",
						})
					}
				}
				return toasts
			},
		},
	)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tpl.Execute(w, data)
	if err != nil {
		fmt.Printf("error executing template: %v", err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}
}

type public interface {
	Public() string
}

type Toast struct {
	Text string
	Type string
}
