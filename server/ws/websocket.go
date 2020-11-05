package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/phamvinhdat/tribe/pkg/wshandler"
	"github.com/sirupsen/logrus"
)

type Websocket interface {
	BroadCastMessage(msg Message)
	RegisterClient(id string, httpResWriter http.ResponseWriter,
		httpReq *http.Request) error
}

type ws struct {
	upgrader websocket.Upgrader
	clients  map[*client]struct{}
}

func New(upgrader websocket.Upgrader) Websocket {
	return &ws{
		upgrader: upgrader,
		clients:  make(map[*client]struct{}),
	}
}

func (h *ws) BroadCastMessage(msg Message) {
	for client, _ := range h.clients {
		err := client.wsHandler.SendMessage(wshandler.Message{
			Timestamp: msg.Timestamp,
			Message:   msg.Message,
		})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":     err,
				"client id": client.id,
			}).Error("failed to send message to client")
		}
	}
}

func (h *ws) RegisterClient(id string, httpResWriter http.ResponseWriter,
	httpReq *http.Request) error {
	conn, err := h.upgrader.Upgrade(httpResWriter, httpReq, nil)
	if err != nil {
		return err
	}

	client := &client{
		id:        id,
		ws:        h,
		conn:      conn,
		wsHandler: wshandler.NewHandler(conn),
	}
	go client.handle()
	h.clients[client] = struct{}{}
	return nil
}
