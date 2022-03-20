package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()

	connectionManager := newConnectionManager()
	go connectionManager.runManager()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!!!")
	})

	app.Post("/internal/notify/:id", func(c *fiber.Ctx) error {
		connectionManager.notify(c.Params("id"), c.Body())
		return c.SendString("Ok")
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		id := c.Params("id")
		connection := newConnection(id, c, func() {
			connectionManager.unregister <- id
		})
		connectionManager.register <- connection
		log.Printf("Registering %s", id)

		go connection.handleIncomingMessages()
		connection.handleMessageSending()
	}))

	log.Fatal(app.Listen(":3002"))
}
