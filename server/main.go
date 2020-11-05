package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/phamvinhdat/tribe/server/router"
	"github.com/phamvinhdat/tribe/server/ws"
)

func main() {
	// new validator
	v := validator.New()

	// websocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(*http.Request) bool { // accept all origin
			return true
		},
	}
	w := ws.New(upgrader)

	// router
	r := gin.Default()
	api := r.Group("/api/v1")
	router.New(w, v).Register(api)

	_ = r.Run()
}
