package main

type ConnectionManager struct {
	connections map[string]*Connection

	register chan *Connection

	unregister chan string
}

func newConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		register:    make(chan *Connection),
		unregister:  make(chan string),
		connections: make(map[string]*Connection),
	}
}

func (connectionManager *ConnectionManager) runManager() {
	for {
		select {
		case connection := <-connectionManager.register:
			connectionManager.connections[connection.id] = connection
		case id := <-connectionManager.unregister:
			if connection, ok := connectionManager.connections[id]; ok {
				close(connection.requestMessageSend)
				delete(connectionManager.connections, connection.id)
			}
		}
	}
}

func (connectionManager *ConnectionManager) notify(id string, msg []byte) {
	if connection, ok := connectionManager.connections[id]; ok {
		connection.requestMessageSend <- msg
	}
}
