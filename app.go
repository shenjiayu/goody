package goody

import (
	"log"
	"net/http"
)

type App struct {
	router *router
	server *Server
}

func NewApp() *App {
	app := new(App)
	app.router = newRouter()
	app.server = newServer()
	return app
}

func (app *App) Handle(pattern string, handler interface{}) error {
	return app.router.Handle(pattern, handler)
}

func (app *App) Run(SSL bool) error {
	http.Handle("/", app.router)
	if err := app.server.Load(); err != nil {
		log.Fatal(err)
	}
	if SSL {
		if err := http.ListenAndServeTLS(app.server.Port, app.server.Cert, app.server.Key, nil); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := http.ListenAndServe(app.server.Port, nil); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
