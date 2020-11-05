package ws

import "time"

type Message struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}
