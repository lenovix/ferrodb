package main

import (
	"log"

	"ferrodb/internal/engine"
	"ferrodb/internal/server"
)

func main() {
	eng := engine.New()

	tcpServer := server.NewTCPServer(":6380", eng)

	if err := tcpServer.Start(); err != nil {
		log.Fatal(err)
	}
}
