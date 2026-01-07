package server

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"ferrodb/internal/config"
	"ferrodb/internal/engine"
)

type TCPServer struct {
	addr     string
	engine   *engine.Engine
	listener net.Listener
	users    []config.User
}

type Client struct {
	conn          net.Conn
	authenticated bool
	user          *config.User
	db            int
}

var rolePermissions = map[string]map[string]bool{
	"admin": {
		"*": true,
	},
	"writer": {
		"SET":     true,
		"DEL":     true,
		"EXPIRE":  true,
		"PERSIST": true,
		"SELECT":  true,
		"INFO":    true,
		"HELP":    true,
	},
	"reader": {
		"GET":  true,
		"TTL":  true,
		"INFO": true,
		"HELP": true,
	},
}

func NewTCPServer(addr string, users []config.User, engine *engine.Engine) *TCPServer {
	return &TCPServer{
		addr:   addr,
		engine: engine,
		users:  users,
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
		db: 0,
	}

	fmt.Fprintln(conn, "Welcome to FerroDB v0.3.4")
	writePrompt(conn, client.db)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			writePrompt(conn, client.db)
			continue
		}

		cmd := strings.Fields(line)
		command := strings.ToUpper(cmd[0])

		if command == "AUTH" {
			if len(cmd) < 3 {
				fmt.Fprintln(conn, "ERR AUTH requires username and password")
				writePrompt(conn, client.db)
				continue
			}

			user := s.findUser(cmd[1], cmd[2])
			if user == nil {
				fmt.Fprintln(conn, "ERR invalid credentials")
				writePrompt(conn, client.db)
				continue
			}

			client.authenticated = true
			client.user = user
			fmt.Fprintf(conn, "OK (%s)\n", user.Role)
			writePrompt(conn, client.db)
			continue
		}

		if !client.authenticated {
			if isPublicCommand(command) {
				result := s.engine.Execute(client.db, line)
				fmt.Fprintln(conn, result)
				writePrompt(conn, client.db)
				continue
			}

			fmt.Fprintln(conn, "NOAUTH Authentication required")
			writePrompt(conn, client.db)
			continue
		}

		if !hasPermission(client.user.Role, command) {
			fmt.Fprintln(conn, "NOPERM permission denied")
			writePrompt(conn, client.db)
			continue
		}

		if command == "SELECT" {
			if len(cmd) < 2 {
				fmt.Fprintln(conn, "ERR SELECT requires index")
				writePrompt(conn, client.db)
				continue
			}

			db, err := strconv.Atoi(cmd[1])
			if err != nil || db < 0 || db > 15 {
				fmt.Fprintln(conn, "ERR invalid DB index")
				writePrompt(conn, client.db)
				continue
			}

			client.db = db
			writePrompt(conn, client.db)
			continue
		}

		if command == "EXIT" {
			fmt.Fprintln(conn, "Bye 👋")
			return
		}

		result := s.engine.Execute(client.db, line)
		fmt.Fprintln(conn, result)
		writePrompt(conn, client.db)
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
	case "HELP", "INFO":
		return true
	default:
		return false
	}
}

func writePrompt(conn net.Conn, db int) {
	fmt.Fprintf(conn, "%d> ", db)
}

func (s *TCPServer) findUser(username, password string) *config.User {
	for i := range s.users {
		u := &s.users[i]
		if u.Username == username && u.Password == password {
			return u
		}
	}
	return nil
}

func hasPermission(role, command string) bool {
	perms := rolePermissions[role]
	if perms["*"] {
		return true
	}
	return perms[command]
}
