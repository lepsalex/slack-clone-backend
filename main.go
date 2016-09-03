/*
 * Main.go
 * Websocket based server implementation
 * for the slack-clone app ...
 *
 * router.go - TODO
 * handlers.go - TODO
 * client.go - TODO
 */

package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type Channel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	router := NewRouter()

	router.Handle("channel add", addChannel)

	http.Handle("/", router)
	http.ListenAndServe(":4000", nil)
}
