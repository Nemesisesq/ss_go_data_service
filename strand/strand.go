package strand

import (
	"net/http"

	"html/template"
)

func ServePage(w http.ResponseWriter, r *http.Request) {
	t := template.New("test template")
	t, _ = t.ParseFiles("strand/linkstation.html", nil)
	t.Execute(w, nil)
}
