package server

import (
	"diyd/src/database/worker"
	"fmt"
	"net/http"
)

type MethodName string

const (
	Get    MethodName = "GET"
	Set    MethodName = "SET"
	Delete MethodName = "DELETE"
)

type Server struct {
	Port   int
	Worker *worker.KVWorker
}

func (s *Server) Start() {
	s.Worker.Start()

	http.HandleFunc("/ping", s.ping)
	http.HandleFunc("/get", s.ping)
	http.HandleFunc("/set", s.ping)
	http.HandleFunc("/delete", s.ping)

	http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil)
}
