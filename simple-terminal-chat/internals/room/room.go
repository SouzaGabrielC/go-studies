package room

import (
	"github.com/google/uuid"
	"net"
	"simple-terminal-chat/internals/client"
	"simple-terminal-chat/internals/message"
	"sync"
)

type Room struct {
	Id        uuid.UUID
	MessageCh chan *message.Message
	clients   *sync.Map
}

func NewRoom() *Room {
	return &Room{
		Id:        uuid.New(),
		MessageCh: make(chan *message.Message, 5),
		clients:   &sync.Map{},
	}
}

func (r *Room) BroadcastToAll(messageContent []byte) {
	r.clients.Range(func(key any, value any) bool {
		conn, ok := value.(net.Conn)
		if !ok {
			return true
		}

		conn.Write(messageContent)

		return true
	})
}

func (r *Room) BroadcastToOthers(senderId uuid.UUID, messageContent []byte) {
	r.clients.Range(func(key any, value any) bool {
		conn, ok := value.(client.ServerClient)
		if !ok {
			return true
		}

		if clientId, ok := key.(uuid.UUID); ok == true && clientId == senderId {
			return true
		}

		conn.Write(messageContent)

		return true
	})
}

func (r *Room) AddClient(client *client.ServerClient) {
	r.clients.Store(client.Id, client)
}

func (r *Room) RemoveClient(clientId uuid.UUID) {
	r.clients.Delete(clientId)
}
