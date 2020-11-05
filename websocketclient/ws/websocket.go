package ws

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/phamvinhdat/tribe/pkg/try"
	"github.com/phamvinhdat/tribe/pkg/wshandler"
	"github.com/sirupsen/logrus"
)

type ws struct {
	dialer    *websocket.Dialer
	conn      *websocket.Conn
	tryer     try.Doer
	url       string
	wsHandler wshandler.Handler

	// unsend message
	unsendMsg *wshandler.Message
}

func New(tryer try.Doer, dialer *websocket.Dialer) *ws {
	return &ws{
		dialer: dialer,
		tryer:  tryer,
	}
}

// Connect
func (w *ws) Connect(url string) error {
	w.url = url
	return w.connect()
}

// Handle Handling when there is an event:
//
// - receive a message
//
// - disconnected
func (w *ws) Handle() {
	w.wsHandler.Handle(
		wshandler.OnReceiveMessage(w.onReceiveMessage),
		wshandler.OnDisConnect(w.onDisconnect),
	)
}

func (w *ws) connect() error {
	conn, _, err := w.dialer.Dial(w.url, nil)
	if err != nil {
		return err
	}

	logrus.Info("connect to websocket server success")
	w.conn = conn
	w.wsHandler = wshandler.NewHandler(conn)
	return nil
}

func (w *ws) tryConnect() {
	err := w.tryer(func() error {
		if err := w.connect(); err != nil {
			logrus.WithField("error", err).
				Error("failed to connect to server, retry ...")
			return try.Continue // retry until timeout
		}

		// re-handle
		go w.Handle()

		// check message unsend
		if w.unsendMsg != nil {
			w.sendMessage(*w.unsendMsg)
			w.unsendMsg = nil // sent
		}

		return nil
	})

	if err != nil {
		logrus.WithField("error", err).Panic("retry failed")
	}
}

// onReceiveMessage
func (w *ws) onReceiveMessage(err error, msg wshandler.Message) {
	if err != nil {
		logrus.WithField("error", err).Error("failed to read message")
		return
	}

	logrus.WithField("message", msg).
		Info("receive message from server")

	// confirm to server
	msg = wshandler.Message{
		Timestamp: time.Now(),
		Message:   "websocket client received message",
	}
	w.sendMessage(msg)
}

// sendMessage send a message via websocket connection
//
// if connection is lost, it will resend when the connection is re-established
func (w *ws) sendMessage(msg wshandler.Message) {
	err := w.wsHandler.SendMessage(msg)
	if err != nil {
		if wshandler.IsLostConnection(err) {
			w.unsendMsg = &msg
			return
		}

		logrus.WithField("error", err).Error("failed to send message")
	}
}

// onDisconnect will try to reconnect until the connection is re-established
func (w *ws) onDisconnect(err error) {
	lg := &logrus.Entry{Logger: logrus.New()}
	if err != nil {
		lg = lg.WithField("error", err)
	}
	lg.Info("disconnected to server, retry ...")
	w.tryConnect()
}
