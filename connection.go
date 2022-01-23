package main

import (
	"encoding/json"

	"github.com/gofiber/websocket/v2"
)

type Connection struct {
	id                 string
	socket             *websocket.Conn
	requestMessageSend chan []byte
	disconnect         func()
}

func newConnection(id string, socket *websocket.Conn, disconnect func()) *Connection {
	return &Connection{
		id:                 id,
		socket:             socket,
		requestMessageSend: make(chan []byte),
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
		if mtype == websocket.BinaryMessage {
			result := struct {
				Relay string `json:"relay"`
			}{}
			json.Unmarshal(msg, &result)
			if result.Relay == "/realtime/status" {
				res, err := json.Marshal(map[string]interface{}{
					"event":  "realtime_status",
					"status": "ok",
				})
				if err == nil {
					connection.requestMessageSend <- res
				}

			} else if result.Relay != "" {
				relay(result.Relay, msg)
			}
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
		connection.socket.WriteMessage(websocket.BinaryMessage, message)

	}
}
