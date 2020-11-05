package wshandler

import "time"

type Message struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

type messageChanModel struct {
	msg Message
	err error
}
