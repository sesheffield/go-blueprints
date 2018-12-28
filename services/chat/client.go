package main

import (
	"github.com/gorilla/websocket"
)

// client represents a single chatting user.
type client struct {
	// socket is the web socket for this client
	socket *websocket.Conn
	// send is a channel on whcih messages are sent
	send chan []byte
	// room is the room this client is chatting in
	room *room
}

// read from the socket, continually sending any received messages to the forward channel on the room type
func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

// write accepts messages from the send chan, writing everything out of the socket
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}
