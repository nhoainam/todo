package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/tuannguyenandpadcojp/fresher26/nam/users/di"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/config"
)

func main() {
	port := getEnv("SERVER_PORT", "50052")

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
