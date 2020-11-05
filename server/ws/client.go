package ws

import (
	"github.com/gorilla/websocket"
	"github.com/phamvinhdat/tribe/pkg/wshandler"
	"github.com/sirupsen/logrus"
)

type client struct {
	id        string
	ws        *ws
	conn      *websocket.Conn
	closeChan chan struct{}
	msgChan   chan Message
	wsHandler wshandler.Handler
}

func (c *client) handle() {
	c.wsHandler.Handle(
		wshandler.OnReceiveMessage(c.onReceiveMessage),
		wshandler.OnDisConnect(c.onDisconnect),
	)
}

func (c *client) onReceiveMessage(err error, msg wshandler.Message) {
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"client id": c.id,
			"error":     err,
		}).Error("failed to read message")
	}

	logrus.WithFields(logrus.Fields{
		"client id": c.id,
		"message":   msg,
	}).Info("receive message from client")
}

func (c *client) onDisconnect(err error) {
	// remove client from websocket hub
	delete(c.ws.clients, c)

	lg := &logrus.Entry{Logger: logrus.New()}
	if err != nil {
		lg = lg.WithField("error", err)
	}
	lg.Infof("client id: %s is disconnected", c.id)
}
