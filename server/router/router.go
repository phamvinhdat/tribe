package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/phamvinhdat/tribe/server/ws"
	"github.com/sirupsen/logrus"
)

type service struct {
	ws        ws.Websocket
	validator *validator.Validate
}

const (
	errorField = "error"
)

func New(ws ws.Websocket, validator *validator.Validate) *service {
	return &service{
		ws:        ws,
		validator: validator,
	}
}

func (s *service) Register(r gin.IRouter) {
	r.POST("/broadcast/msg", s.postBroadcastMsg)
	r.GET("/ws", s.getWebSocketConnection)
}

func (s *service) getWebSocketConnection(c *gin.Context) {
	clientID := uuid.New().String()
	if err := s.ws.RegisterClient(clientID, c.Writer, c.Request); err != nil {
		logrus.WithField(errorField, err).
			Error("failed to upgrade connection")
		return
	}
}

func (s *service) postBroadcastMsg(c *gin.Context) {
	// bind message
	var msg MessageBroadcast
	if err := c.ShouldBind(&msg); err != nil {
		logrus.WithField(errorField, err).
			Error("failed to bind broadcast message")
		c.JSON(http.StatusBadRequest, gin.H{
			errorField: "failed to bind message",
		})
		return
	}

	logrus.WithField("content", msg).
		Info("receive message from publishing client")

	// validate message
	if err := s.validator.Struct(msg); err != nil {
		logrus.WithField(errorField, err).
			Error("message invalid")
		c.JSON(http.StatusBadRequest, gin.H{
			errorField: "message invalid",
		})
		return
	}

	// broadcast message
	s.ws.BroadCastMessage(ws.Message{
		Timestamp: msg.Timestamp,
		Message:   msg.Message,
	})

	// return http status
	c.Status(http.StatusNoContent)
}
