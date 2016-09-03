/*
 * CLIENT.GO
 * Defines interface for client (API for server)
 * and is primarily responsible for seding messages
 * to and from the browser
 */

package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Client - Defines channel struct
type Client struct {
	// socket *websocket.Conn
	send chan Message
}

func (client *Client) write() {
	for msg := range client.send {
		// TODO socket.sendJSON
		fmt.Printf("%#v\n", msg)
	}
}

func (client *Client) subscribeChannels() {
	// TODO change feed Query rethinkDB
	for {
		time.Sleep(r())
		client.send <- Message{"channel add", ""}
	}
}

func (client *Client) subscribeMessages() {
	// TODO change feed Query rethinkDB
	for {
		time.Sleep(r())
		client.send <- Message{"message add", ""}
	}
}

// TEMP
func r() time.Duration {
	return time.Millisecond * time.Duration(rand.Intn(1000))
}

// NewClient - Function that creates object (similar to constructor but no)
func NewClient() *Client {
	return &Client{
		send: make(chan Message),
	}
}
