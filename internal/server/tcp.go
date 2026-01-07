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

	"golang.org/x/crypto/bcrypt"
)

type TCPServer struct {
	addr     string
	engine   *engine.Engine
	listener net.Listener
	users    []config.User
	dbCount  int
}

type Client struct {
	conn          net.Conn
	authenticated bool
	user          *config.User
	db            int
	resp          bool
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

func NewTCPServer(
	addr string,
	users []config.User,
	dbCount int,
	engine *engine.Engine,
) *TCPServer {
	return &TCPServer{
		addr:    addr,
		users:   users,
		dbCount: dbCount,
		engine:  engine,
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

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimSpace(line)

		if line == "" {
			writePrompt(conn, client.db)
			continue
		}

		if line[0] == '*' {
			args, err := readRESPFromLine(reader, line)
			if err != nil {
				writeError(conn, "ERR invalid RESP")
				continue
			}
			s.handleRESP(conn, client, args)
			continue
		}

		s.handleInline(conn, client, line)
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
	case "HELP", "INFO", "EXIT":
		return true
	default:
		return false
	}
}

func writePrompt(conn net.Conn, db int) {
	fmt.Fprintf(conn, "%d> ", db)
}

func (s *TCPServer) findUser(username, password string) *config.User {
	for _, u := range s.users {
		if u.Username != username {
			continue
		}

		err := bcrypt.CompareHashAndPassword(
			[]byte(u.Password),
			[]byte(password),
		)
		if err == nil {
			return &u
		}
	}
	return nil
}

func hasPermission(role, command string) bool {
	switch command {
	case "AUTH", "HELP", "INFO", "ACL", "LOGOUT":
		return true
	}

	perms := rolePermissions[role]
	if perms["*"] {
		return true
	}
	return perms[command]
}

// func (s *TCPServer) handleCommandInline(
// 	conn net.Conn,
// 	client *Client,
// 	line string,
// ) {
// 	result := s.execute(client, strings.Fields(line))
// 	fmt.Fprintln(conn, result)
// 	writePrompt(conn, client.db)
// }

// func (s *TCPServer) handleCommandRESP(
// 	conn net.Conn,
// 	client *Client,
// 	args []string,
// ) {
// 	result, kind := s.executeRESP(client, args)

// 	switch kind {
// 	case "ok":
// 		writeSimpleString(conn, result)
// 	case "err":
// 		writeError(conn, result)
// 	case "int":
// 		n, _ := strconv.ParseInt(result, 10, 64)
// 		writeInteger(conn, n)
// 	case "bulk":
// 		if result == "" {
// 			writeNull(conn)
// 		} else {
// 			writeBulkString(conn, result)
// 		}
// 	}
// }

// func (s *TCPServer) executeRESP(client *Client, args []string) (string, string) {
// 	cmd := strings.ToUpper(args[0])

// 	result := s.engine.Execute(client.db, strings.Join(args, " "))

// 	switch cmd {
// 	case "GET":
// 		if result == "(nil)" {
// 			return "", "bulk"
// 		}
// 		return result, "bulk"

// 	case "TTL":
// 		return result, "int"

// 	default:
// 		if strings.HasPrefix(result, "ERR") {
// 			return result, "err"
// 		}
// 		return result, "ok"
// 	}
// }

func (s *TCPServer) execute(
	client *Client,
	args []string,
) (string, string) {

	cmd := strings.ToUpper(args[0])

	// ===== AUTH =====
	if cmd == "AUTH" {
		if client.authenticated {
			return "ERR already authenticated (use LOGOUT)", "err"
		}
		if len(args) < 3 {
			return "ERR AUTH username password", "err"
		}
		user := s.findUser(args[1], args[2])
		if user == nil {
			return "ERR invalid credentials", "err"
		}
		client.authenticated = true
		client.user = user
		return "OK", "ok"
	}

	// ===== PUBLIC =====
	if !client.authenticated && !isPublicCommand(cmd) {
		return "NOAUTH Authentication required", "err"
	}

	// ===== PERMISSION =====
	if client.authenticated && !hasPermission(client.user.Role, cmd) {
		return "NOPERM permission denied", "err"
	}

	// ===== LOGOUT =====
	if cmd == "LOGOUT" {
		client.authenticated = false
		client.user = nil
		client.db = 0
		return "OK logged out", "ok"
	}

	// ===== QUIT / EXIT =====
	if cmd == "QUIT" || cmd == "EXIT" {
		return "BYE", "close"
	}

	// ===== SELECT =====
	if cmd == "SELECT" {
		if len(args) < 2 {
			return "ERR SELECT index", "err"
		}
		db, err := strconv.Atoi(args[1])
		if err != nil || db < 0 || db >= s.dbCount {
			return "ERR invalid DB index", "err"
		}
		client.db = db
		return "OK", "ok"
	}

	// ===== ACL =====
	if cmd == "ACL" {
		if len(args) < 2 {
			return "ERR ACL subcommand required", "err"
		}
		sub := strings.ToUpper(args[1])

		switch sub {
		case "WHOAMI":
			return fmt.Sprintf(
				"user=%s role=%s",
				client.user.Username,
				client.user.Role,
			), "bulk"

		case "CAT":
			perms := rolePermissions[client.user.Role]
			var out []string
			for p := range perms {
				out = append(out, p)
			}
			return strings.Join(out, " "), "bulk"

		default:
			return "ERR unknown ACL subcommand", "err"
		}
	}

	// ===== ENGINE =====
	res := s.engine.Execute(client.db, strings.Join(args, " "))

	if strings.HasPrefix(res, "ERR") {
		return res, "err"
	}

	// RESP-aware type
	switch cmd {
	case "GET":
		if res == "(nil)" {
			return "", "null"
		}
		return res, "bulk"
	case "TTL":
		return res, "int"
	default:
		return res, "ok"
	}
}

func (s *TCPServer) handleRESP(conn net.Conn, client *Client, args []string) {
	result, kind := s.execute(client, args)

	switch kind {
	case "ok":
		writeSimpleString(conn, result)
	case "err":
		writeError(conn, result)
	case "int":
		n, _ := strconv.ParseInt(result, 10, 64)
		writeInteger(conn, n)
	case "bulk":
		writeBulkString(conn, result)
	case "null":
		writeNull(conn)
	case "close":
		writeSimpleString(conn, result)
		conn.Close()
	}
}

func (s *TCPServer) handleInline(conn net.Conn, client *Client, line string) {
	args := strings.Fields(line)
	if len(args) == 0 {
		writePrompt(conn, client.db)
		return
	}

	result, kind := s.execute(client, args)

	if kind == "close" {
		fmt.Fprintln(conn, result)
		conn.Close()
		return
	}

	fmt.Fprintln(conn, result)
	writePrompt(conn, client.db)
}

func readRESPFromLine(
	r *bufio.Reader,
	firstLine string,
) ([]string, error) {

	// firstLine contoh: "*3"
	if len(firstLine) < 2 || firstLine[0] != '*' {
		return nil, fmt.Errorf("invalid RESP array")
	}

	count, err := strconv.Atoi(firstLine[1:])
	if err != nil || count < 0 {
		return nil, fmt.Errorf("invalid RESP array length")
	}

	args := make([]string, 0, count)

	for i := 0; i < count; i++ {
		// expect: $<len>
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)

		if len(line) < 2 || line[0] != '$' {
			return nil, fmt.Errorf("invalid RESP bulk string")
		}

		size, err := strconv.Atoi(line[1:])
		if err != nil || size < 0 {
			return nil, fmt.Errorf("invalid bulk size")
		}

		// read exact <size> bytes + \r\n
		buf := make([]byte, size+2)
		if _, err := r.Read(buf); err != nil {
			return nil, err
		}

		arg := string(buf[:size])
		args = append(args, arg)
	}

	return args, nil
}
