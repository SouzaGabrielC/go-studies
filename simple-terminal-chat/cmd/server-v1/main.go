package main

import (
	"fmt"
	"github.com/google/uuid"
	"io"
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
	slog.Info(fmt.Sprintf("[%s] Connection established. Num of connections", connId))
	defer func(conn net.Conn) {
		connections.Delete(connId)
		err := conn.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("[%s] Close connection error. Error: %s", connId, err.Error()))
			return
		}

		slog.Info(fmt.Sprintf("[%s] Connection closed.", connId))
	}(conn)

	messageContentBytes := make([]byte, 128)

	var pr *io.PipeReader
	var pw *io.PipeWriter

	for {
		n, err := conn.Read(messageContentBytes)
		if err != nil {
			if err != io.EOF {
				slog.Error(fmt.Sprintf("[%s] Error on receiving the messageContentBytes: %s", connId, err.Error()))
			}
			break
		}

		if n == 0 {
			continue
		}

		if pw == nil {
			pr, pw = io.Pipe()
			messagePipeCh <- &message.Message{
				MessageReader: pr,
				SenderId:      connId,
			}
		}

		pw.Write(messageContentBytes[:n])

		if messageContentBytes[n-1] == '\n' {
			pw.Close()
			pw = nil
			pr = nil
		}
	}
}

func broadcastMessageToConnections(connections *sync.Map, messagePipeCh <-chan *message.Message) {
	for {
		receivedMessage := <-messagePipeCh
		messageContent := make([]byte, 128)

		for {
			n, err := receivedMessage.MessageReader.Read(messageContent)

			if err != nil {
				break
			}

			connections.Range(func(key any, value any) bool {
				conn, ok := value.(net.Conn)
				if !ok {
					return true
				}

				if connId, ok := key.(uuid.UUID); ok == true && connId == receivedMessage.SenderId {
					return true
				}

				conn.Write(messageContent[:n])

				return true
			})
		}

	}
}
