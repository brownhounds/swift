package res

import (
	"html/template"
	"net/http"
)

func HtmlTemplate(w http.ResponseWriter, status int, path string, data any) {
	tpl, err := template.ParseFiles(path)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	if err := tpl.Execute(w, data); err != nil {
		panic(err)
	}
}
