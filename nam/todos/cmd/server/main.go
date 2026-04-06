package main

// main.go — Application Entry Point
//
// This file is responsible for:
// 1. Loading configuration from environment variables (envconfig)
// 2. Initializing all dependencies via Wire (see di/wire.go)
// 3. Setting up the gRPC server with interceptors
// 4. Starting the server and handling graceful shutdown
//
// You will build this incrementally:
//   - Phase 1: Understand the entry point pattern
//   - Phase 1: Add gRPC server setup and interceptor chain
//   - Phase 2: Add Wire-based dependency injection (replace manual wiring)
//   - Phase 5: Add observability (Zap logger, Datadog tracer, Sentry)
//
// See: resources/phase-01-architecture-grpc.md (gRPC server setup, clean architecture)
// See: resources/phase-02-database-di.md (Wire DI — replaces the manual wiring below)

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/di"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/config"
)

func main() {
	port := getEnv("SERVER_PORT", "50051")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcSrv, cleanup, err := di.InitializeServer(cfg)
	if err != nil {
		log.Fatalf("failed to create gRPC server: %v", err)
	}
	defer cleanup()

	// --- Graceful shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("gRPC server listening on :%s", port)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	<-quit
	log.Println("shutdown signal received, stopping gRPC server...")
	grpcSrv.GracefulStop()
	log.Println("server stopped")
}

// getEnv returns the value of the environment variable named by key,
// or defaultVal if the variable is not set.
func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
