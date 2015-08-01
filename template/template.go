package template

import (
	"html/template"
	"net/http"
)

var funcs_map = template.FuncMap{
	"backtohtml": backtohtml,
	"multiply":   multiply,
}

func backtohtml(data string) interface{} {
	return template.HTML(data)
}

func multiply(a, b float32) float32 {
	return a * b
}

func RenderTemplate(w http.ResponseWriter, page string, data interface{}) {
	t := template.New("")
	t.Funcs(funcs_map)
	t.ParseFiles("view/header.html", "view/footer.html", "view/"+page+".html")
	if err := t.ExecuteTemplate(w, page+".html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
