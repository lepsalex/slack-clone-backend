/*
 * CLIENT.GO
 * Defines interface for client (API for server)
 * and is primarily responsible for seding messages
 * to and from the browser
 */

package main

import (
	"log"

	"github.com/gorilla/websocket"
	r "gopkg.in/dancannon/gorethink.v2"
)

// FindHandler function signature definition
type FindHandler func(string) (Handler, bool)

// Client - Defines channel struct
type Client struct {
	send         chan Message
	socket       *websocket.Conn
	findHandler  FindHandler
	session      *r.Session
	stopChannels map[int]chan bool
	id           string
	userName     string
}

// NewStopChannel - Stop channel creator
func (client *Client) NewStopChannel(stopKey int) chan bool {
	// Ensure we stop and existing stop channel on this key
	client.StopForKey(stopKey)

	// Start stop channel on key
	stop := make(chan bool)
	client.stopChannels[stopKey] = stop
	return stop
}

// StopForKey - stops channel goroutine on key
func (client *Client) StopForKey(key int) {
	if ch, found := client.stopChannels[key]; found {
		ch <- true
		delete(client.stopChannels, key)
	}
}

// Message - Defines message structure
type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

// Read method (client)
func (client *Client) Read() {
	var message Message
	for {
		if err := client.socket.ReadJSON(&message); err != nil {
			break
		}
		if handler, found := client.findHandler(message.Name); found == true {
			handler(client, message.Data)
		}
	}
	client.socket.Close()
}

// Write method (client)
func (client *Client) Write() {
	for msg := range client.send {
		if err := client.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	client.socket.Close()
}

// Close method (client)
func (client *Client) Close() {
	for _, ch := range client.stopChannels {
		ch <- true
	}
	close(client.send)
}

// NewClient - Function that creates object (similar to constructor but no)
func NewClient(socket *websocket.Conn, findHanlder FindHandler, session *r.Session) *Client {
	// Default client user values
	var user User
	user.Name = "anonymous"
	// Insert default into rethinkDB
	res, err := r.Table("user").Insert(user).RunWrite(session)
	if err != nil {
		log.Println(err.Error())
	}
	// Get generated user key
	var id string
	if len(res.GeneratedKeys) > 0 {
		id = res.GeneratedKeys[0]
	}
	// Return client for user
	return &Client{
		send:         make(chan Message),
		socket:       socket,
		findHandler:  findHanlder,
		session:      session,
		stopChannels: make(map[int]chan bool),
		id:           id,
		userName:     user.Name,
	}
}
