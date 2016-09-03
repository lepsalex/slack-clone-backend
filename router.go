/*
 * ROUTER.GO
 * Takes incoming websocket http
 * and routes to appropriate handler
 * function (handlers.go)
 */

package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Handler - Function signature definition
type Handler func(*Client, interface{})

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Router - Router struct
type Router struct {
	rules map[string]Handler
}

// NewRouter - Init function for Router Object
func NewRouter() *Router {
	return &Router{
		rules: make(map[string]Handler),
	}
}

// Handle - Assigns handler to router msg mapping
func (r *Router) Handle(msgName string, handler Handler) {
	r.rules[msgName] = handler
}

func (e *Router) ServeHTTP(w http.ResponseWrite, r *http.Request) {
	// Upgrade all traffic from HTTP to Websocket
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// handle errors if any
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	// Init new client
	client := NewClient(socket)
	go client.Write()
	client.Read()
}
