package msgservice

import "time"

type Message struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type MsgServerRes struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}
