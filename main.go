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

import "net/http"

// Message - Defines message structure
type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

// Channel - Defines channel structure
type Channel struct {
	ID   string `json:"ID"`
	Name string `json:"name"`
}

func main() {
	router := NewRouter()

	router.Handle("channel add", addChannel)

	http.Handle("/", router)
	http.ListenAndServe(":4000", nil)
}
