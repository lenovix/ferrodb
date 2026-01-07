package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"ferrodb/internal/config"
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

	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	eng := engine.New(cfg)

	// 🔴 TCP (redis-cli)
	tcpServer := server.NewTCPServer(
		cfg.Server.Address,
		cfg.Users,
		cfg.Engine.DBCount,
		eng,
	)

	// 🟢 HTTP (Web UI / API)
	httpServer := server.NewHTTPServer(
		":8080",
		eng,
	)

	go func() {
		if err := tcpServer.Start(); err != nil {
			log.Println("TCP server error:", err)
			stop()
		}
	}()

	go func() {
		httpServer.Start()
	}()

	<-ctx.Done()
	log.Println("🛑 Shutting down FerroDB...")

	tcpServer.Shutdown()
	eng.Shutdown()

	log.Println("✅ FerroDB stopped gracefully")
}
