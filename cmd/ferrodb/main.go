package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"ferrodb/internal/adminapi"
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

	// Load config
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	// Init engine
	eng := engine.New(cfg)

	// 🔴 TCP Server (RESP / redis-cli)
	tcpServer := server.NewTCPServer(
		cfg.Server.Address,
		cfg.Users,
		cfg.Engine.DBCount,
		eng,
	)

	// 🟢 Admin HTTP API (Web UI backend)
	adminapi.SetEngine(eng)
	adminServer := adminapi.New(":8080")
	adminServer.Start()

	// Start TCP server
	go func() {
		if err := tcpServer.Start(); err != nil {
			log.Println("TCP server error:", err)
			stop()
		}
	}()

	<-ctx.Done()
	log.Println("🛑 Shutting down FerroDB...")

	tcpServer.Shutdown()
	eng.Shutdown()

	log.Println("✅ FerroDB stopped gracefully")
}
