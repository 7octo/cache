package routes

import (
    "threat-detection/internal/handler"
    
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App, eventHandler *handler.EventHandler) {
    // Middleware
    app.Use(logger.New())
    
    // Routes
    api := app.Group("/api/v1")
    api.Get("/health", eventHandler.HealthCheck)
    api.Post("/events", eventHandler.AddEvent)
}