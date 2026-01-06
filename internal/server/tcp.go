package server

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"ferrodb/internal/engine"
)

type TCPServer struct {
	addr     string
	engine   *engine.Engine
	listener net.Listener
}

func NewTCPServer(addr string, engine *engine.Engine) *TCPServer {
	return &TCPServer{
		addr:   addr,
		engine: engine,
	}
}

func (s *TCPServer) Start() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.listener = ln
	log.Println("🚀 FerroDB TCP server running on", s.addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			log.Println("accept error:", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Fprintln(conn, "Welcome to FerroDB v0.2.1")
	fmt.Fprintln(conn, "Type HELP for commands")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.ToUpper(line) == "EXIT" {
			fmt.Fprintln(conn, "Bye 👋")
			return
		}

		result := s.engine.Execute(line)
		fmt.Fprintln(conn, result)
	}
}

func (s *TCPServer) Shutdown() {
	if s.listener != nil {
		log.Println("🔌 Closing TCP listener")
		s.listener.Close()
	}
}
