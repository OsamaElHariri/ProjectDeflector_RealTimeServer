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

	app.Use("/", func(c *fiber.Ctx) error {
		userId := c.Get("x-user-id")
		if userId != "" {
			c.Locals("userId", userId)
		}
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!!!")
	})

	app.Post("/internal/notify/:id", func(c *fiber.Ctx) error {
		connectionManager.notify(c.Params("id"), c.Body())
		return c.SendString("Ok")
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		playerId := c.Locals("userId").(string)
		connection := newConnection(playerId, c, func() {
			connectionManager.unregister <- playerId
		})
		connectionManager.register <- connection
		log.Printf("Registering %s", playerId)

		go connection.handleIncomingMessages()
		connection.handleMessageSending()
	}))

	log.Fatal(app.Listen(":3002"))
}
