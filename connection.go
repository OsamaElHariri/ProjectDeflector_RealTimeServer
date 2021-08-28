package main

import (
	"github.com/gofiber/websocket/v2"
)

type Connection struct {
	id                 string
	socket             *websocket.Conn
	requestMessageSend chan string
	disconnect         func()
}

func newConnection(id string, socket *websocket.Conn, disconnect func()) *Connection {
	return &Connection{
		id:                 id,
		socket:             socket,
		requestMessageSend: make(chan string),
		disconnect:         disconnect,
	}
}

func (connection *Connection) handleIncomingMessages() {
	for {
		mtype, msg, err := connection.socket.ReadMessage()
		if err != nil {
			connection.disconnect()
			break
		}
		if mtype == websocket.TextMessage {
			connection.socket.WriteMessage(mtype, msg)
		}
	}
}

func (connection *Connection) handleMessageSending() {
	for {
		message, ok := <-connection.requestMessageSend
		if !ok {
			connection.socket.WriteMessage(websocket.CloseMessage, []byte{})
			connection.socket.Close()
			return
		}
		connection.socket.WriteMessage(websocket.TextMessage, []byte(message))

	}
}
