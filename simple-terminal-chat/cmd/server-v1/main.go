package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net"
	"simple-terminal-chat/internals/message"
	"sync"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	messagePipeCh := make(chan *message.Message, 5)

	connections := &sync.Map{}

	go broadcastMessageToConnections(connections, messagePipeCh)

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("Connection error. Error:", err.Error())
		}

		connId := uuid.New()
		connections.Store(connId, conn)

		go handleConnectionRead(connId, conn, messagePipeCh, connections)
	}
}

func handleConnectionRead(connId uuid.UUID, conn net.Conn, messagePipeCh chan<- *message.Message, connections *sync.Map) {
	slog.Info(fmt.Sprintf("[%s] Connection established.", connId))
	defer func(conn net.Conn) {
		connections.Delete(connId)
		err := conn.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("[%s] Close connection error. Error: %s", connId, err.Error()))
			return
		}

		slog.Info(fmt.Sprintf("[%s] Connection closed.", connId))
	}(conn)

	connScanner := bufio.NewScanner(conn)

	for connScanner.Scan() {
		var msgContent message.Content
		err := json.Unmarshal(connScanner.Bytes(), &msgContent)
		if err != nil {
			slog.Error("Invalid message received. Error:", err.Error())
			continue
		}

		messagePipeCh <- &message.Message{
			SenderId:       connId,
			MessageContent: msgContent,
		}
	}
}

func broadcastMessageToConnections(connections *sync.Map, messagePipeCh <-chan *message.Message) {
	for {
		receivedMessage := <-messagePipeCh

		connections.Range(func(key any, value any) bool {
			conn, ok := value.(net.Conn)
			if !ok {
				return true
			}

			if connId, ok := key.(uuid.UUID); ok == true && connId == receivedMessage.SenderId {
				return true
			}

			json.NewEncoder(conn).Encode(receivedMessage.MessageContent)

			return true
		})

	}
}
