/*
 * HANDLER.GO
 * Defines all handlers mapped to
 * expected message names
 */

package main

import (
	"github.com/mitchellh/mapstructure"
	r "gopkg.in/dancannon/gorethink.v2"
)

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
