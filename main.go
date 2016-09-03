/*
 * MAIN.GO
 * Websocket based server implementation
 * for the slack-clone app ...
 *
 * router.go - TODO
 * handlers.go - TODO
 * client.go - TODO
 */

package main

import "net/http"

// Channel - Defines channel structure
type Channel struct {
	ID   string `json:"ID"`
	Name string `json:"name"`
}

func main() {
	router := NewRouter()

	// Register route handlers
	router.Handle("channel add", addChannel)

	// Activate router
	http.Handle("/", router)
	http.ListenAndServe(":4000", nil)
}
