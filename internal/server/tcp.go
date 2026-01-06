package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"ferrodb/internal/engine"
)

type TCPServer struct {
	addr   string
	engine *engine.Engine
}

func NewTCPServer(addr string, engine *engine.Engine) *TCPServer {
	return &TCPServer{
		addr:   addr,
		engine: engine,
	}
}

func (s *TCPServer) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	fmt.Println("🚀 FerroDB TCP server running on port", s.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
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
