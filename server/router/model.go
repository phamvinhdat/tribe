package router

import "time"

// MessageBroadcast is the data required to 'broadcast' the message
type MessageBroadcast struct {
	Message   string    `json:"message" validate:"required"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
}

type HTTPResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}
