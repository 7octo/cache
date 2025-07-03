package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "threat-detection/internal/config"
    "threat-detection/internal/detector"
    "threat-detection/internal/handler"
    "threat-detection/internal/model"
    "threat-detection/routes"
    
    "github.com/gofiber/fiber/v2"
)

// ExampleThreatChecker implements the ThreatChecker interface
type ExampleThreatChecker struct{}

func (c *ExampleThreatChecker) IsThreat(event interface{}) bool {
    e, ok := event.(model.Event)
    if !ok {
        return false
    }
    
    // Example threat detection - check for "malicious" in payload
    if payload, ok := e.Payload["data"].(string); ok {
        return payload == "malicious"
    }
    return false
}

func (c *ExampleThreatChecker) HandleThreat(event interface{}) {
    e, ok := event.(model.Event)
    if !ok {
        return
    }
    
    log.Printf("THREAT DETECTED - Event ID: %s, Source: %s", e.ID, e.Metadata.Source)
    // In a real app, you might notify security team, block IP, etc.
}

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Create threat detector
    threatChecker := &ExampleThreatChecker{}
    detector := detector.NewCachedThreatDetector(threatChecker)
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