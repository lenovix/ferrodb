package server

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"ferrodb/internal/engine"
)

type TCPServer struct {
	addr     string
	engine   *engine.Engine
	listener net.Listener
	password string
}

type Client struct {
	conn          net.Conn
	authenticated bool
	db            int
}

func NewTCPServer(addr, password string, engine *engine.Engine) *TCPServer {
	return &TCPServer{
		addr:     addr,
		engine:   engine,
		password: password,
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

	client := &Client{
		conn:          conn,
		authenticated: s.password == "",
		db:            0,
	}

	fmt.Fprintln(conn, "Welcome to FerroDB v0.3.2")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		cmd := strings.Fields(line)
		command := strings.ToUpper(cmd[0])

		// 🔐 AUTH
		if !client.authenticated && !isPublicCommand(command) {
			fmt.Fprintln(conn, "NOAUTH Authentication required")
			continue
		}

		if command == "AUTH" {
			if len(cmd) < 2 {
				fmt.Fprintln(conn, "ERR AUTH requires password")
				continue
			}

			if cmd[1] == s.password {
				client.authenticated = true
				fmt.Fprintln(conn, "OK")
			} else {
				fmt.Fprintln(conn, "ERR invalid password")
			}
			continue
		}

		// 🗂 SELECT DB
		if command == "SELECT" {
			if len(cmd) < 2 {
				fmt.Fprintln(conn, "ERR SELECT requires index")
				continue
			}

			db, err := strconv.Atoi(cmd[1])
			if err != nil || db < 0 || db > 15 {
				fmt.Fprintln(conn, "ERR invalid DB index")
				continue
			}

			client.db = db
			fmt.Fprintln(conn, "OK")
			continue
		}

		// 👋 EXIT
		if command == "EXIT" {
			fmt.Fprintln(conn, "Bye 👋")
			return
		}

		// ⚙ ENGINE EXECUTION
		result := s.engine.Execute(client.db, line)
		fmt.Fprintln(conn, result)
	}
}

func (s *TCPServer) Shutdown() {
	if s.listener != nil {
		log.Println("🔌 Closing TCP listener")
		s.listener.Close()
	}
}

func isPublicCommand(cmd string) bool {
	switch cmd {
	case "AUTH", "HELP", "INFO":
		return true
	default:
		return false
	}
}
