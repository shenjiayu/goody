package goody

import (
	"encoding/json"
	"os"
)

type Server struct {
	Port string `json:"port"`
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

func newServer() *Server {
	return &Server{Port: ":8080"}
}

func (this *Server) Load() error {
	reader, err := os.Open("server.json")
	buf := make([]byte, 1024)
	if err != nil {
		return err
	}
	n, err := reader.Read(buf)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf[:n], this)
	if err != nil {
		return err
	}
	return nil
}
