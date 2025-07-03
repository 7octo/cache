package model

type Event struct {
    ID      string                 `json:"id"`
    Payload map[string]interface{} `json:"payload"`
    Metadata struct {
        Timestamp time.Time `json:"timestamp"`
        Source    string    `json:"source"`
    } `json:"metadata"`
}