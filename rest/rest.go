package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"chainspace.io/blockmania/node"
	"chainspace.io/blockmania/rest/api"
	"chainspace.io/blockmania/rest/service"
)

type Server struct {
	port   uint
	router *api.Router
	srv    *http.Server
	wg     *sync.WaitGroup
}

func New(port uint, node *node.Server) *Server {
	srv := service.New(node)
	router := api.New(srv)
	httpsrv := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: router,
	}
	return &Server{port, router, httpsrv, &sync.WaitGroup{}}
}

func (s *Server) Start() {
	s.wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		log.Printf("http server started on port %v", s.port)
		log.Printf("http server exited: %v", s.srv.ListenAndServe())
	}(s.wg)
}

func (s *Server) Shutdown() error {
	err := s.srv.Shutdown(context.Background())
	s.wg.Wait()
	return err
}