/*
 * CLIENT.GO
 * Defines interface for client (API for server)
 * and is primarily responsible for seding messages
 * to and from the browser
 */

package main

import "github.com/gorilla/websocket"

// FindHandler function signature definition
type FindHandler func(string) (Handler, bool)

// Client - Defines channel struct
type Client struct {
	send        chan Message
	socket      *websocket.Conn
	findHandler FindHandler
}

// Message - Defines message structure
type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

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

func (client *Client) Write() {
	for msg := range client.send {
		if err := client.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	client.socket.Close()
}

// NewClient - Function that creates object (similar to constructor but no)
func NewClient(socket *websocket.Conn, findHanlder FindHandler) *Client {
	return &Client{
		send:        make(chan Message),
		socket:      socket,
		findHandler: findHanlder,
	}
}
