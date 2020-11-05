package main

import (
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/phamvinhdat/tribe/pkg/try"
	"github.com/phamvinhdat/tribe/websocketclient/ws"
	"github.com/sirupsen/logrus"
)

func main() {
	dialer := websocket.DefaultDialer
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:8080",
		Path:   "/api/v1/ws",
	}

	// tryer
	tryer := try.New(
		try.WithInterval(30*time.Second), // retry every 30s
		try.WithTimeout(24*time.Hour),    // limited timeout to avoid endless loops (1 day)
	)

	// websocket handler
	webs := ws.New(tryer, dialer)
	err := webs.Connect(u.String())
	if err != nil {
		logrus.WithField("error", err).
			Panic("failed to connect to websocket server")
	}

	go webs.Handle()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	close(quit)
	logrus.Info("shutting down")
}
