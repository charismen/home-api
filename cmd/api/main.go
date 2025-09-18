package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charismen/home-api/internal/handlers"
	"github.com/charismen/home-api/internal/repository"
	"github.com/charismen/home-api/internal/service"
	"github.com/charismen/home-api/pkg/apiclient"
	"github.com/charismen/home-api/pkg/database"
	"github.com/charismen/home-api/pkg/redis"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Initialize database
	database.InitDB()
	defer database.CloseDB()

	// Initialize Redis
	redis.InitRedis()
	defer redis.CloseRedis()

	// Create API client
	apiClient := apiclient.NewClient("https://pokeapi.co/api/v2")

	// Create repository
	itemRepo := repository.NewItemRepository()

	// Create service
	itemService := service.NewItemService(itemRepo, apiClient)

	// Create handler
	itemHandler := handlers.NewItemHandler(itemService)

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/sync", itemHandler.SyncItems).Methods("POST")
	r.HandleFunc("/items", itemHandler.GetItems).Methods("GET")

	// Set up background job
	c := cron.New()
	_, err = c.AddFunc("*/15 * * * *", func() {
		log.Println("Running scheduled sync job")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		err := itemService.SyncItems(ctx)
		if err != nil {
			log.Printf("Scheduled sync failed: %v", err)
		} else {
			log.Println("Scheduled sync completed successfully")
		}
	})
	if err != nil {
		log.Fatalf("Failed to set up cron job: %v", err)
	}
	c.Start()
	defer c.Stop()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Server shutting down...")
	
	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Server exited properly")
}