package detector

import "time"

type ThreatChecker interface {
    IsThreat(event interface{}) bool
    HandleThreat(event interface{})
}

type EventDetector interface {
    AddEvent(eventID string, event interface{})
    Start()
    Stop()
}