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
	r "gopkg.in/dancannon/gorethink.v2"
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
	rules   map[string]Handler
	session *r.Session
}

// NewRouter - Init function for Router Object
func NewRouter(session *r.Session) *Router {
	return &Router{
		rules:   make(map[string]Handler),
		session: session,
	}
}

// Handle - Assigns handler to router msg mapping
func (router *Router) Handle(msgName string, handler Handler) {
	router.rules[msgName] = handler
}

// FindHandler - function to assist in look up of router handlers
func (router *Router) FindHandler(msgName string) (Handler, bool) {
	handler, found := router.rules[msgName]
	return handler, found
}

// ServeHTTP - define our implementation of the ServerHTTP func
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Upgrade all traffic from HTTP to Websocket
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// handle errors if any
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	// Init new client
	client := NewClient(socket, router.FindHandler, router.session)
	go client.Write()
	client.Read()
}
