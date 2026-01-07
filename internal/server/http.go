package server

import (
	"encoding/json"
	"log"
	"net/http"
)

type HTTPServer struct {
	addr   string
	engine EngineAdapter
}

type EngineAdapter interface {
	DBCount() int
}

func NewHTTPServer(addr string, engine EngineAdapter) *HTTPServer {
	return &HTTPServer{
		addr:   addr,
		engine: engine,
	}
}

func (s *HTTPServer) Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"version": "0.5.0",
		})
	})

	log.Println("🌐 HTTP server running on", s.addr)
	http.ListenAndServe(s.addr, mux)
}
