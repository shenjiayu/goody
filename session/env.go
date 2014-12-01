package session

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type Env struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Session        *Session
	Output_method  string
	Output_data    interface{}
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

func (e *Env) SetHeader(w http.ResponseWriter, key string, value string) {
	header := w.Header()
	header.Set(key, value)
}

func (e *Env) NotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

//---------part to be removed to template folder
var funcs_map = template.FuncMap{
	"backtohtml": backtohtml,
}

func backtohtml(data string) interface{} {
	return template.HTML(data)
}

func (e *Env) RenderTemplate(w http.ResponseWriter, page string, data interface{}) {
	t := template.New("")
	t.Funcs(funcs_map)
	t.ParseFiles("view/header.html", "view/footer.html", "view/"+page+".html")
	if err := t.ExecuteTemplate(w, page+".html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//

//Output_method is used to determine the way to response to the client, such as render a template or return 'json' message.
//render represents for rendering the template
func (e *Env) Set_output_method(method string, data interface{}) {
	e.Output_method = method
	e.Output_data = data
}

func (e *Env) ServeJson(v interface{}, w http.ResponseWriter) {
	output, _ := json.Marshal(v)
	fmt.Fprintf(w, "%s", output)
}

func (e *Env) ServeEventStream(v interface{}, w http.ResponseWriter) {
	e.SetHeader(w, "Content-Type", "text/event-stream")
	e.SetHeader(w, "Cache-Control", "no-cache")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	output, _ := json.Marshal(v)
	fmt.Fprintf(w, "data: %s\n\n", output)
	flusher.Flush()
}

func Encrypt(data interface{}) string {
	h := sha1.New()
	buf := make([]byte, 5)
	switch data.(type) {
	case int64:
		binary.PutVarint(buf, data.(int64))
	case string:
		buf = []byte(data.(string))
	}
	h.Write(buf)
	return fmt.Sprintf("%x", h.Sum(nil))
}
