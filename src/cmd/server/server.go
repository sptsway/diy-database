package server

import (
	"diyd/src/cmd/app"
	"diyd/src/config"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type MethodName string

const (
	Get    MethodName = "GET"
	Set    MethodName = "SET"
	Delete MethodName = "DELETE"
)

type Server struct {
	cmdcfg  config.CmdArgs
	Port    int
	closers []io.Closer
}

func (s *Server) Start() {
	fmt.Println("Server starting on port: ", s.Port)

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	kvapp, err := app.NewApp(s.cmdcfg)
	if err != nil {
		panic(err)
		return
	}
	if closer, ok := kvapp.(io.Closer); ok {
		s.closers = append(s.closers, closer)
	}

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: mux,
	}

	mux.HandleFunc("/ping", kvapp.Ping)
	mux.HandleFunc("/get", kvapp.Get)
	mux.HandleFunc("/set", kvapp.Set)
	mux.HandleFunc("/delete", kvapp.Delete)

	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			s.Stop()
			panic(err)
			return
		}
	}()

	<-stop
	s.Stop()
}

func (s *Server) Stop() {
	fmt.Println("Server shutting down")

	for _, closer := range s.closers {
		err := closer.Close()
		if err != nil {
			fmt.Println("Failed to Close closer", err)
		}
	}

	fmt.Println("Server shut down")
}
