package goody

import (
	"fmt"
	"net/http"
)

type App struct {
	router *router
}

func NewApp() *App {
	app := new(App)
	app.router = newRouter()
	return app
}

func (app *App) Handle(pattern string, handler interface{}) error {
	return app.router.Handle(pattern, handler)
}

func (app *App) ProcessRequest(function interface{}) {
	app.router.processRequest(function)
}

func (app *App) Run(SSL bool) error {
	http.Handle("/", app.router)
	if SSL {
		if err := http.ListenAndServeTLS(":8080", "/etc/nginx/ssl/coddict.crt", "/etc/nginx/ssl/server.key", nil); err != nil {
			fmt.Println(err)
		}
	} else {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
