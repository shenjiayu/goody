package goody

import (
	"testing"
)

type mainController struct {
}

func (this *mainController) Get() {

}

func (this *mainController) Post() {

}

func TestHandle(t *testing.T) {
	app := NewApp()
	app.Handle("/", new(mainController))
}
