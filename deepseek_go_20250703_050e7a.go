package detector

import (
    "sync"
    "time"
    
    "github.com/patrickmn/go-cache"
)

type CachedThreatDetector struct {
    cache           *cache.Cache
    checkIntervals  []time.Duration
    threatChecker   ThreatChecker
    cleanupInterval time.Duration
    mu              sync.Mutex
    stopChan        chan struct{}
}

type cachedEvent struct {
    Event      interface{}
    NextCheck  time.Time
    CheckIndex int
}

func NewCachedThreatDetector(checker ThreatChecker) *CachedThreatDetector {
    intervals := []time.Duration{
        1 * time.Minute,
        2 * time.Minute,
        3 * time.Minute,
        5 * time.Minute,
        10 * time.Minute,
    }
    
    return &CachedThreatDetector{
        cache:           cache.New(24*time.Hour, 30*time.Minute),
        checkIntervals:  intervals,
        threatChecker:   checker,
        cleanupInterval: 24 * time.Hour,
        stopChan:        make(chan struct{}),
    }
}

func (d *CachedThreatDetector) AddEvent(eventID string, event interface{}) {
    d.mu.Lock()
    defer d.mu.Unlock()
    
    d.cache.Set(eventID, &cachedEvent{
        Event:      event,
        NextCheck:  time.Now().Add(d.checkIntervals[0]),
        CheckIndex: 0,
    }, cache.DefaultExpiration)
}

func (d *CachedThreatDetector) Start() {
    go d.runChecker()
}

func (d *CachedThreatDetector) Stop() {
    close(d.stopChan)
}

func (d *CachedThreatDetector) runChecker() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            d.checkEvents()
        case <-d.stopChan:
            return
        }
    }
}

func (d *CachedThreatDetector) checkEvents() {
    d.mu.Lock()
    defer d.mu.Unlock()
    
    now := time.Now()
    
    for _, item := range d.cache.Items() {
        eventID := item.Key
        cachedEvt, ok := item.Object.(*cachedEvent)
        if !ok {
            d.cache.Delete(eventID)
            continue
        }
        
        if now.After(cachedEvt.NextCheck) {
            if d.threatChecker.IsThreat(cachedEvt.Event) {
                d.threatChecker.HandleThreat(cachedEvt.Event)
                d.cache.Delete(eventID)
            } else {
                nextIndex := cachedEvt.CheckIndex + 1
                if nextIndex < len(d.checkIntervals) {
                    cachedEvt.NextCheck = now.Add(d.checkIntervals[nextIndex])
                    cachedEvt.CheckIndex = nextIndex
                    d.cache.Set(eventID, cachedEvt, cache.DefaultExpiration)
                } else {
                    d.cache.Delete(eventID)
                }
            }
        }
    }
}