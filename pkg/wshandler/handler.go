package wshandler

import (
	"errors"

	"github.com/gorilla/websocket"
)

type (
	OnReceiveMessageFn func(err error, msg Message)
	OnDisconnectFn     func(err error)
	OnSendMessageFn    func(msg Message) error

	Handler interface {
		Handle(opts ...Option)
		SendMessage(msg Message) error
		Close()
	}

	wsHandler struct {
		conn      *websocket.Conn
		closeChan chan error
		msgChan   chan messageChanModel
		isClose   bool
	}
)

var lostConnection = errors.New("connection is lost")

func IsLostConnection(err error) bool {
	return lostConnection == err
}

func NewHandler(conn *websocket.Conn) Handler {
	return &wsHandler{
		conn:      conn,
		closeChan: make(chan error),
		msgChan:   make(chan messageChanModel),
	}
}

// SendMessage send a message and notify if connection is lost
func (h *wsHandler) SendMessage(msg Message) error {
	err := h.conn.WriteJSON(msg)
	if err != nil {
		// check for lost connection
		if websocket.IsUnexpectedCloseError(err) {
			h.closeChan <- err
			return lostConnection
		}

		return err
	}

	return nil
}

// Handle
func (h *wsHandler) Handle(opts ...Option) {
	// get option
	opt := option{
		onReMsgFn: func(error, Message) {}, // default do nothing
		onDisFn:   func(error) {},          // default do nothing
	}
	for _, o := range opts {
		o.apply(&opt)
	}

	h.onDisconnect()
	h.onMessage()
	for !h.isClose {
		select {
		case err := <-h.closeChan:
			opt.onDisFn(err)
		case msgChanModel := <-h.msgChan:
			opt.onReMsgFn(msgChanModel.err, msgChanModel.msg)
		}
	}
}

func (h *wsHandler) Close() {
	if h.isClose {
		return
	}

	close(h.closeChan)
	close(h.msgChan)
	h.isClose = true
}

func (h *wsHandler) onMessage() {
	// receive message
	go func() {
		for {
			var msg Message
			err := h.conn.ReadJSON(&msg)
			if err != nil {
				// check for lost connection
				if websocket.IsUnexpectedCloseError(err) {
					h.closeChan <- err
					return
				}
			}

			h.msgChan <- messageChanModel{
				msg: msg,
				err: err,
			}
		}
	}()
}

// onDisconnect handle for client send close signal
func (h *wsHandler) onDisconnect() {
	h.conn.SetCloseHandler(func(int, string) error {
		h.closeChan <- nil
		return nil
	})
}
