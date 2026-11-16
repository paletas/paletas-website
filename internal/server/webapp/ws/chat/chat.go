package chat

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	"github.com/paletas/paletas_website/internal/concurrent"
)

var connections = concurrent.NewConcurrentSlice()
var roomsOnline []*ChatRoom

func ConfigureChat(app *fiber.App) {

	timer := time.NewTicker((time.Second * 10))

	go func() {
		for range timer.C {
			for _, channel := range roomsOnline {
				if len(channel.Users) > 0 {
					channel.broadcastChannelStatus()
				}
			}
		}
	}()

	app.Use("/api/chat/:channel?", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/api/chat/:channel?", websocket.New(func(conn *websocket.Conn) {
		// Handle WebSocket connections
		log.Println("WebSocket connection established")

		username := conn.Query("username")
		user := NewChatUser(username, conn)

		channelName := conn.Params("channel")
		if channelName != "" {
			channel := findChannel(channelName)
			if channel == nil {
				channel = NewChatRoom(channelName, user)
				roomsOnline = append(roomsOnline, channel)
			} else {
				channel.Join(user)
				user.sendChannelStatus(channel)
			}
		}

		connections.Append(conn)

		for {
			messageType, message, error := conn.ReadMessage()
			log.Println("WebSocket message received")

			if error != nil {
				log.Println("WebSocket read error:", error)
				break
			}

			if messageType == websocket.TextMessage {
				chatOperation, deserializeError := UnmarshalOperation(message)
				if deserializeError != nil {
					log.Println("WebSocket deserialize error:", deserializeError)
					break
				}

				chatOperation.Perform(user)
				log.Println("WebSocket message broadcasted")
			}
		}

		// Close the WebSocket connection
		log.Println("WebSocket connection closed")
		connections.RemoveItem(conn)
	}))
}

func findChannel(name string) *ChatRoom {
	for _, room := range roomsOnline {
		if room.Name == name {
			return room
		}
	}
	return nil
}
