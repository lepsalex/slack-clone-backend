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

import (
	"net/http"

	"log"

	r "gopkg.in/dancannon/gorethink.v2"
)

// Channel - Defines channel structure
type Channel struct {
	ID   string `json:"ID" gorethink:"id,omitempty"`
	Name string `json:"name" gorethink:"name"`
}

// User - defines user type
type User struct {
	ID   string `gorethink:"id,omitempty"`
	Name string `gorethink:"name"`
}

func main() {
	// Connect to rethinkDB
	session, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "slack_clone",
	})

	if err != nil {
		log.Panic(err.Error())
	}

	// Create new router object
	router := NewRouter(session)

	// Register route handlers
	router.Handle("channel add", addChannel)
	router.Handle("channel subscribe", subscribeChannel)

	// Activate router
	http.Handle("/", router)
	http.ListenAndServe(":4000", nil)
}
