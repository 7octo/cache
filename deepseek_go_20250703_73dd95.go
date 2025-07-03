package handler

import (
    "threat-detection/internal/detector"
    "threat-detection/internal/model"
    
    "github.com/gofiber/fiber/v2"
)

type EventHandler struct {
    detector detector.EventDetector
}

func NewEventHandler(d detector.EventDetector) *EventHandler {
    return &EventHandler{detector: d}
}

func (h *EventHandler) AddEvent(c *fiber.Ctx) error {
    var event model.Event
    
    if err := c.BodyParser(&event); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Cannot parse JSON",
        })
    }
    
    if event.ID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Event ID is required",
        })
    }
    
    h.detector.AddEvent(event.ID, event)
    
    return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
        "message": "Event received",
        "event_id": event.ID,
    })
}

func (h *EventHandler) HealthCheck(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status": "ok",
    })
}