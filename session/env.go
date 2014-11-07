package session

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

var (
	t = template.Must(template.ParseFiles("view/header.html",
		"view/footer.html",
		"view/home.html",
		"view/register.html",
		"view/login.html",
		"view/ads.html",
		"view/all_ads.html",
		"view/new_ads.html",
		"view/all_events.html",
		"view/admin_dashboard.html",
		"view/profile.html",
	))
)

type Env struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Session        *Session
	Csrf_required  bool
	Tpl            string
	Output         interface{}
}

func NewEnv(w http.ResponseWriter, r *http.Request) *Env {
	env := new(Env)
	env.ResponseWriter = w
	env.Request = r
	return env
}

func (e *Env) Redirect(url string) {
	http.Redirect(e.ResponseWriter, e.Request, url, http.StatusFound)
}

func (e *Env) Set_csrf(required bool) {
	e.Csrf_required = required
}

func (e *Env) Set_tpl(tpl string) {
	e.Tpl = tpl
}

func (e *Env) RenderTemplate(w http.ResponseWriter, page string, data interface{}) {
	if err := t.ExecuteTemplate(w, page+".html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (e *Env) Set_output(body interface{}) {
	e.Output = body
}

func (e *Env) OutputJson(v interface{}, w http.ResponseWriter) {
	output, _ := json.Marshal(v)
	fmt.Fprintf(w, "%s", output)
}
