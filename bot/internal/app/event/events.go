package event

import "encoding/json"

type EventType string

// Обрабатываемые типы событий
const (
	EventMessageNew EventType = "message_new"
)

type GroupEvent struct {
	EventID string          `json:"event_id"`
	GroupID int             `json:"group_id"`
	Object  json.RawMessage `json:"object"`
	Type    EventType       `json:"type"`
	V       string          `json:"v"`
}
