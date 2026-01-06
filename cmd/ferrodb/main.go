package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"ferrodb/internal/engine"
	"ferrodb/internal/server"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	eng := engine.New()
	tcpServer := server.NewTCPServer(":6380", eng)

	go func() {
		if err := tcpServer.Start(); err != nil {
			log.Println("server error:", err)
			stop()
		}
	}()

	<-ctx.Done()
	log.Println("🛑 Shutting down FerroDB...")

	tcpServer.Shutdown()
	eng.Shutdown()

	log.Println("✅ FerroDB stopped gracefully")
}
