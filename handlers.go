/*
 * HANDLER.GO
 * Defines all handlers mapped to
 * expected message names
 */

package main

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	r "gopkg.in/dancannon/gorethink.v2"
)

// Stop channel key
const (
	ChannelStop = iota
	UserStop
	MessageStop
)

// Adds a channel to the app
func addChannel(client *Client, data interface{}) {
	var channel Channel
	err := mapstructure.Decode(data, &channel)
	if err != nil {
		client.send <- Message{"error", err.Error()}
		return
	}
	go func() {
		err = r.Table("channel").
			Insert(channel).
			Exec(client.session)
		if err != nil {
			client.send <- Message{"error", err.Error()}
		}
	}()
}

// Subscribes client (user) to changes on channel
func subscribeChannel(client *Client, data interface{}) {
	stop := client.NewStopChannel(ChannelStop)
	result := make(chan r.ChangeResponse)
	cursor, err := r.Table("channel").
		Changes(r.ChangesOpts{IncludeInitial: true}).
		Run(client.session)
	if err != nil {
		client.send <- Message{"error", err.Error()}
	}
	go func() {
		var change r.ChangeResponse
		for cursor.Next(&change) {
			result <- change
		}
	}()
	go func() {
		for {
			select {
			case <-stop:
				// Echo stop
				cursor.Close()
				return
			case change := <-result:
				if change.NewValue != nil && change.OldValue == nil {
					client.send <- Message{"channel add", change.NewValue}
					fmt.Println("Sent Channel Add")
				}
			}
		}
	}()
}

// Unsubscribes client (user) to changes on channel
func unsubscribeChannel(client *Client, data interface{}) {
	client.StopForKey(ChannelStop)
}

/*
 * TODO
 * user edit
 * user subscribe
 * user unsubscribe
 * message add
 * message subscribe
 * message unsubscribe
 */
