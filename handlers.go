/*
 * HANDLER.GO
 * Defines all handlers mapped to
 * expected message names
 */

package main

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	r "gopkg.in/dancannon/gorethink.v2"
)

// Stop channel keys
const (
	ChannelStop = iota
	UserStop
	MessageStop
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

// ChannelMessage - defines message type
type ChannelMessage struct {
	ID        string    `gorethink:"id,omitempty"`
	ChannelID string    `gorethink:"channelID"`
	Body      string    `gorethink:"body"`
	Author    string    `gorethink:"author"`
	CreatedAt time.Time `gorethink:"createdAt"`
}

/*
 * CHANNEL HANDLERS
 */

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

// Unsubscribes client (user) from changes on channel
func unsubscribeChannel(client *Client, data interface{}) {
	client.StopForKey(ChannelStop)
}

/*
 * USER HANDLERS
 */

// Edit client user name
func editUser(client *Client, data interface{}) {
	var user User
	err := mapstructure.Decode(data, &user)
	if err != nil {
		client.send <- Message{"error", err.Error()}
		return
	}
	client.userName = user.Name
	go func() {
		_, err := r.Table("users").
			Get(client.id).
			Update(user).
			RunWrite(client.session)
		if err != nil {
			client.send <- Message{"error", err.Error()}
		}
	}()
}

// Subscribe client (user) to changes on User
func subscribeUser(client *Client, data interface{}) {
	go func() {
		stop := client.NewStopChannel(UserStop)
		cursor, err := r.Table("user").
			Changes(r.ChangesOpts{IncludeInitial: true}).
			Run(client.session)
		if err != nil {
			client.send <- Message{"error", err.Error()}
			return
		}
		changeFeedHelper(cursor, "user", client.send, stop)
	}()
}

// Unsubscribes client (user) from changes on user
func unsubscribeUser(client *Client, data interface{}) {
	client.StopForKey(UserStop)
}

/*
 * MESSAGE HANDLERS
 */

// Adds a message to a channel
func addChannelMessage(client *Client, data interface{}) {
	var channelMessage ChannelMessage
	err := mapstructure.Decode(data, &channelMessage)
	if err != nil {
		client.send <- Message{"error", err.Error()}
	}
	go func() {
		channelMessage.CreatedAt = time.Now()
		channelMessage.Author = client.userName
		err = r.Table("messages").
			Insert(channelMessage).
			Exec(client.session)
		if err != nil {
			client.send <- Message{"error", err.Error()}
		}
	}()
}

// Subscribe client (user) to message on a channel
func subscribeChannelMessage(client *Client, data interface{}) {
	go func() {
		eventData := data.(map[string]interface{})
		val, ok := eventData["channelId"]
		if !ok {
			return
		}
		channelID, ok := val.(string)
		if !ok {
			return
		}
		stop := client.NewStopChannel(MessageStop)
		cursor, err := r.Table("messages").
			OrderBy(r.OrderByOpts{Index: r.Desc("createdAt")}).
			Filter(r.Row.Field("channelId").Eq(channelID)).
			Changes(r.ChangesOpts{IncludeInitial: true}).
			Run(client.session)
		if err != nil {
			client.send <- Message{"error", err.Error()}
			return
		}
		changeFeedHelper(cursor, "message", client.send, stop)
	}()
}

// Unsubscribes client (user) from changes on messages
func unsubscribeChannelMessage(client *Client, data interface{}) {
	client.StopForKey(MessageStop)
}

/*
 * HELPERS
 */

// Change feed helper function
func changeFeedHelper(cursor *r.Cursor, changeEventName string, send chan<- Message, stop <-chan bool) {
	change := make(chan r.ChangeResponse)
	cursor.Listen(change)
	for {
		eventName := ""
		var data interface{}
		select {
		// If we receive stop channel then stop
		case <-stop:
			cursor.Close()
			return
		// Otherwise figure out the type of action by comparing new vs. old values from rethinkDB
		case val := <-change:
			if val.NewValue != nil && val.OldValue == nil {
				// No old value but a new value = add
				eventName = changeEventName + " add"
				data = val.NewValue
			} else if val.NewValue == nil && val.OldValue != nil {
				// No new value but an old value = remove
				eventName = changeEventName + " remove"
				data = val.OldValue
			} else if val.NewValue != nil && val.OldValue != nil {
				// Both old and new values = edit
				eventName = changeEventName + " edit"
				data = val.NewValue
			}
			send <- Message{eventName, data}
		}
	}
}

/*
 * TODO
 * message add
 * message subscribe
 * message unsubscribe
 * refactor channel to match user and the rest
 */
