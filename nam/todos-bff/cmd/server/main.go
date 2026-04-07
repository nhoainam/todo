package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/di"
)

// main.go — BFF Application Entry Point
// Phase 3: GraphQL & BFF Pattern

func main() {
	// Initialize dependencies via Wire
	httpServer, cleanup, err := di.InitializeServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
	defer cleanup()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Mount the GraphQL handler
	router.Handle("/query", httpServer)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("BFF GraphQL Server started on port %s", port)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down BFF server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
