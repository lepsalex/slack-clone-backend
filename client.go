/*
 * CLIENT.GO
 * Defines interface for client (API for server)
 * and is primarily responsible for seding messages
 * to and from the browser
 */

package main

import (
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
}

// NewStopChannel - Stop channel creator
func (client *Client) NewStopChannel(stopKey int) chan bool {
	stop := make(chan bool)
	client.stopChannels[stopKey] = stop
	return stop
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
func NewClient(socket *websocket.Conn,
	findHanlder FindHandler,
	session *r.Session) *Client {
	return &Client{
		send:         make(chan Message),
		socket:       socket,
		findHandler:  findHanlder,
		session:      session,
		stopChannels: make(map[int]chan bool),
	}
}
