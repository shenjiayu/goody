package goody

import (
	"encoding/json"
	"fmt"
	"github.com/shenjiayu/goody/session"
	"github.com/shenjiayu/goody/template"
	"net/http"
)

type Controller struct {
}

func (this *Controller) Redirect(w http.ResponseWriter, r *http.Request, url string, code int) {
	http.Redirect(w, r, url, code)
}

func (this *Controller) SetHeader(w http.ResponseWriter, key, value string) {
	w.Header().Set(key, value)
}

func (this *Controller) NotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (this *Controller) serveJson(w http.ResponseWriter, data interface{}) {
	output, _ := json.Marshal(data)
	fmt.Fprintf(w, "%s", output)
}

func (this *Controller) Respond(w http.ResponseWriter, method string, data interface{}, ctx session.Context) {
	switch method {
	case "render":
		template.RenderTemplate(w, data.(string), ctx)
	case "json":
		this.serveJson(w, data)
	default:
		fmt.Errorf("No such methods, only support [render, json]")
	}
}
