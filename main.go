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
	router.Handle("channel unsubscribe", unsubscribeChannel)

	router.Handle("user edit", editUser)
	router.Handle("user subscribe", subscribeUser)
	router.Handle("user unsubscribe", unsubscribeUser)

	router.Handle("message add", addChannelMessage)
	router.Handle("message subscribe", subscribeChannelMessage)
	router.Handle("message unsubscribe", unsubscribeChannelMessage)

	// Activate router
	http.Handle("/", router)
	http.ListenAndServe(":4000", nil)
}
