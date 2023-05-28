package chat

import (
	"log"

	"github.com/gofiber/websocket/v2"
)

type ChatRoom struct {
	Name  string      `json:"name"`
	Users []*ChatUser `json:"users"`
}

type ChatUser struct {
	Username   string `json:"username"`
	Connection *websocket.Conn
}

func NewChatRoom(name string, user *ChatUser) *ChatRoom {
	return &ChatRoom{
		Name:  name,
		Users: []*ChatUser{user},
	}
}

func NewChatUser(username string, conn *websocket.Conn) *ChatUser {
	return &ChatUser{
		Username:   username,
		Connection: conn,
	}
}

func (room *ChatRoom) Join(user *ChatUser) {
	room.Users = append(room.Users, user)
}

func (room *ChatRoom) Leave(user *ChatUser) {
	for i, u := range room.Users {
		if u.Username == user.Username {
			room.Users = append(room.Users[:i], room.Users[i+1:]...)
			break
		}
	}
}

func (user *ChatUser) sendChannelStatus(channel *ChatRoom) {
	var currentStatus = NewChannelStatusUpdateOperation(channel.Name, len(channel.Users))

	err := user.Connection.WriteJSON(currentStatus)
	if err != nil {
		log.Println("WebSocket write error:", err)
	}
}

func (channel *ChatRoom) broadcastChannelStatus() {
	var currentStatus = NewChannelStatusUpdateOperation(channel.Name, len(channel.Users))

	for _, user := range channel.Users {
		err := user.Connection.WriteJSON(currentStatus)
		if err != nil {
			log.Println("WebSocket write error:", err)
		}
	}
}

func (channel *ChatRoom) broadcastMessage(message string) {
	var newMessage = NewBroadcastMessageOperation(channel.Name, message)

	for _, user := range channel.Users {
		err := user.Connection.WriteJSON(newMessage)
		if err != nil {
			log.Println("WebSocket write error:", err)
		}
	}
}
