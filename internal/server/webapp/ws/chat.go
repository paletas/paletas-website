package ws

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	"github.com/paletas/paletas_website/internal/concurrent"
)

const (
	NewMessageOperationType         = "newMessage"
	ServerStatusUpdateOperationType = "serverStatusUpdate"
)

type ChatOperation struct {
	Operation string `json:"operation"`
}

type NewMessageOperation struct {
	ChatOperation
	Message string `json:"message"`
}

type ServerStatusUpdateOperation struct {
	ChatOperation
	UsersOnline int `json:"usersOnline"`
}

var connections *concurrent.ConcurrentSlice

func ConfigureChat(app *fiber.App) {

	connections = concurrent.NewConcurrentSlice()

	timer := time.NewTicker((time.Second * 10))

	go func() {
		for range timer.C {
			if connections.Len() > 0 {
				broadcastServerStatus()
			}
		}
	}()

	app.Use("/api/chat", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/api/chat", websocket.New(func(conn *websocket.Conn) {
		// Handle WebSocket connections
		log.Println("WebSocket connection established")

		connections.Append(conn)
		sendServerStatus(conn)

		for {
			messageType, message, error := conn.ReadMessage()
			log.Println("WebSocket message received")

			if error != nil {
				log.Println("WebSocket read error:", error)
				break
			}

			if messageType == websocket.TextMessage {
				broadcastMessage(string(message))
				log.Println("WebSocket message broadcasted")
			}
		}

		// Close the WebSocket connection
		log.Println("WebSocket connection closed")
		connections.RemoveItem(conn)
	}))
}

func sendServerStatus(conn *websocket.Conn) {
	var currentStatus = &ServerStatusUpdateOperation{
		ChatOperation: ChatOperation{Operation: ServerStatusUpdateOperationType},
		UsersOnline:   connections.Len(),
	}

	err := conn.WriteJSON(currentStatus)
	if err != nil {
		log.Println("WebSocket write error:", err)
	}
}

func broadcastServerStatus() {
	var currentStatus = &ServerStatusUpdateOperation{
		ChatOperation: ChatOperation{Operation: ServerStatusUpdateOperationType},
		UsersOnline:   connections.Len(),
	}

	for conn := range connections.Iter() {
		err := conn.Value.(*websocket.Conn).WriteJSON(currentStatus)
		if err != nil {
			log.Println("WebSocket write error:", err)
		}
	}
}

func broadcastMessage(message string) {
	var newMessage = &NewMessageOperation{
		ChatOperation: ChatOperation{Operation: NewMessageOperationType},
		Message:       message,
	}

	for conn := range connections.Iter() {
		err := conn.Value.(*websocket.Conn).WriteJSON(newMessage)
		if err != nil {
			log.Println("WebSocket write error:", err)
		}
	}
}
