package context

import (
	"encoding/json"
	"github.com/coddict/session"
	"net/http"
	"strconv"
)

type Env struct {
	Request   *http.Request
	Status    int
	Ctx       Context
	Session   *session.Session
	ResWriter http.ResponseWriter
	finished  bool
}

func NewEnv(w http.ResponseWriter, r *http.Request) *Env {
	e := new(Env)
	e.Request = r
	e.ResWriter = w
	e.Ctx = NewContext()
	e.finished = false
	e.Status = http.StatusOK
	return e
}

func (e *Env) Header() http.Header {
	return e.ResWriter.Header()
}

func (e *Env) SetContentType(tp string) {
	e.ResWriter.Header().Set("Content-type", tp)
}

func (e *Env) SetContentJson() {
	e.ResWriter.Header().Set("Content-type", "application/json; charset=utf-8")
}

func (e *Env) SetStatus(status int) {
	e.Status = status
}

func (e *Env) Write(v interface{}) {
	if e.finished {
		return
	}
	buf, err := json.Marshal(v)
	if err != nil {
		e.WriteError(http.StatusInternalServerError, err)
	} else {
		e.SetContentJson()
		e.write(buf)
	}
}

func (e *Env) WriteString(data string) {
	if len(e.Header().Get("Content-type")) == 0 {
		e.SetContentType("text/plain")
	}
	e.write([]byte(data))
}

func (e *Env) WriteBuffer(data []byte) {
	if len(e.Header().Get("Content-type")) == 0 {
		e.SetContentType("application/octet-stream")
	}
	e.write(data)
}

func (e *Env) write(data []byte) {
	if e.finished {
		return
	}
	e.finished = true
	e.ResWriter.Header().Set("Content-Length", strconv.Itoa(len(data)))
	e.ResWriter.WriteHeader(e.Status)
	e.ResWriter.Write(data)
}

func (e *Env) WriteError(status int, err error) {
	e.Status = status
	e.WriteString(err.Error())
}

func (e *Env) Redirect(url string, status int) {
	e.finished = true
	http.Redirect(e.ResWriter, e.Request, url, status)
}

func (e *Env) SetCookie(c *http.Cookie) {
	http.SetCookie(e.ResWriter, c)
}

func (e *Env) Finish() {
	if e.finished {
		return
	}
	e.finished = true
	e.ResWriter.WriteHeader(e.Status)
}

func (e *Env) IsFinished() bool {
	return e.finished
}
