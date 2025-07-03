package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	
	"threat-detection/internal/config"
	"threat-detection/internal/detector"
	"threat-detection/internal/handler"
	"threat-detection/routes"
	
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	// Create threat service and processor
	threatService := detector.NewHTTPThreatService("https://threat-api.example.com/check")
	threatProcessor := detector.NewThreatProcessor(threatService)
	
	// Create threat detector
	detector := detector.NewCachedThreatDetector(threatProcessor)
	detector.Start()
	defer detector.Stop()
	
	// Create Fiber app
	app := fiber.New()
	
	// Create handlers
	eventHandler := handler.NewEventHandler(detector)
	
	// Setup routes
	routes.SetupRoutes(app, eventHandler)
	
	// Start server
	go func() {
		if err := app.Listen(cfg.Server.Port); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}
