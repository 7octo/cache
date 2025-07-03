package detector

import (
	"context"
	"time"
)

type Threat struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Score     float64                `json:"score"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type ThreatChecker interface {
	CheckForThreats(ctx context.Context, event interface{}) ([]Threat, error)
	ProcessThreat(ctx context.Context, threat Threat) error
}

type EventDetector interface {
	AddEvent(eventID string, event interface{})
	Start()
	Stop()
}
