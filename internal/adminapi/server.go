package adminapi

import (
	"log"
	"net/http"
)

type Server struct {
	addr string
}

func New(addr string) *Server {
	return &Server{addr: addr}
}

func (s *Server) Start() {
	mux := http.NewServeMux()
	RegisterRoutes(mux)

	log.Println("[admin] listening on", s.addr)
	go http.ListenAndServe(s.addr, mux)
}
